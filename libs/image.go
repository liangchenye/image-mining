package libs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	clairHost                  = "http://127.0.0.1:6060"
	postLayerURI               = "/v1/layers"
	getLayerVulnerabilitiesURI = "/v1/layers/%s/vulnerabilities?minimumPriority=%s"
	httpPort                   = 9279
)

type APIVulnerabilitiesResponse struct {
	Vulnerabilities []APIVulnerability
}

type APIVulnerability struct {
	ID, Link, Priority, Description string
}

type Layer struct {
	ID          string
	ParentID    string
	Path        string
	ImageFormat string
}

//Only support 'docker' at present
type Image struct {
	Format string //Docker or ACI
	User   string
	Repo   string
	Tag    string
	ID     string
	Path   string //the saved dir
	Layers []string
}

func ImageNew(format string, user string, repo string, tag string) (Image, error) {
	if format == "" {
		return Image{}, errors.New("'Format' is mandatory.")
	}

	if repo == "" {
		return Image{}, errors.New("'Repo' is mandatory.")
	}

	if tag == "" {
		return Image{}, errors.New("'Tag' is mandatory.")
	}
	image := Image{Format: format, User: user, Repo: repo, Tag: tag}

	return image, nil
}

func (image *Image) Pull() error {
	if cached, err := image.cached(); err != nil {
		return err
	} else if cached {
		return nil
	}
	_, err := ExecCmd("", "docker", "pull", fmt.Sprintf("%s:%s", image.Repo, image.Tag))
	return err
}

func (image *Image) Scan() error {
	if _, err := image.Save(); err != nil {
		return err
	}

	if _, err := image.History(); err != nil {
		return err
	}

	for i := 0; i < len(image.Layers); i++ {
		var err error
		if i > 0 {
			err = analyzeLayer(clairHost, image.Path+"/"+image.Layers[i]+"/layer.tar", image.Layers[i], image.Layers[i-1], "Docker")
		} else {
			err = analyzeLayer(clairHost, image.Path+"/"+image.Layers[i]+"/layer.tar", image.Layers[i], "", "Docker")
		}
		if err != nil {
			fmt.Println("- Could not analyze layer: %s\n", err)
		}
	}
	return nil
}

func (image *Image) GetVuln() {
	vulnerabilities, err := getVulnerabilities(clairHost, image.ID, "Low")
	if err != nil {
		fmt.Println("- Could not get vulnerabilities: %s\n", err)
	}
	if len(vulnerabilities) == 0 {
		fmt.Println("Bravo, your image looks SAFE !")
	}
	for _, vulnerability := range vulnerabilities {
		fmt.Printf("- # %s\n", vulnerability.ID)
		fmt.Printf("  - Priority:    %s\n", vulnerability.Priority)
		fmt.Printf("  - Link:        %s\n", vulnerability.Link)
		fmt.Printf("  - Description: %s\n", vulnerability.Description)
	}

}

func (image *Image) Save() (string, error) {
	imageName := image.GetID()
	if imageName == "" {
		fmt.Println("Cannot find the image, try to pull.")
		return "", errors.New("Cannot find the image, try to pull.")
	}
	//My data volume.
	image.Path = path.Join("/image-data/docker-images", "official", imageName)
	topTar := path.Join(image.Path, imageName, "layer.tar")
	if _, err := os.Stat(topTar); err == nil {
		fmt.Println("Already saved")
		return image.Path, nil
	}

	os.MkdirAll(image.Path, 0777)

	//Code from github.com/coreos/clair/contrib/analyze-local-images
	var stderr bytes.Buffer
	save := exec.Command("docker", "save", imageName)
	save.Stderr = &stderr
	extract := exec.Command("tar", "xf", "-", "-C"+image.Path)
	extract.Stderr = &stderr
	pipe, err := extract.StdinPipe()
	if err != nil {
		return "", err
	}
	save.Stdout = pipe

	err = extract.Start()
	if err != nil {
		return "", errors.New(stderr.String())
	}
	err = save.Run()
	if err != nil {
		return "", errors.New(stderr.String())
	}
	err = pipe.Close()
	if err != nil {
		return "", err
	}
	err = extract.Wait()
	if err != nil {
		return "", errors.New(stderr.String())
	}

	return image.Path, nil
}

func (image *Image) History() ([]string, error) {
	imageName := image.GetID()
	if imageName == "" {
		fmt.Println("Cannot find the image, try to pull.")
		return nil, errors.New("Cannot find the image, try to pull.")
	}
	var stderr bytes.Buffer
	cmd := exec.Command("docker", "history", "-q", "--no-trunc", imageName)
	cmd.Stderr = &stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return []string{}, err
	}

	err = cmd.Start()
	if err != nil {
		return []string{}, errors.New(stderr.String())
	}

	var layers []string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		layers = append(layers, scanner.Text())
	}

	for i := len(layers)/2 - 1; i >= 0; i-- {
		opp := len(layers) - 1 - i
		layers[i], layers[opp] = layers[opp], layers[i]
	}

	image.Layers = layers
	return layers, nil
}

func (image *Image) Clear() {
}

func (image *Image) GetID() string {
	var repo string
	if image.ID != "" {
		return image.ID
	}
	if image.User == "" {
		repo = image.Repo
	} else {
		repo = fmt.Sprintf("%s/%s", image.User, image.Repo)
	}
	if out, err := ExecCmd("", "docker", "images", "--no-trunc=true", repo); err == nil {
		lines := strings.Split(out, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, repo) {
				fields := strings.Fields(line)
				if len(fields) < 3 {
					continue
				} else if fields[1] == image.Tag {
					image.ID = fields[2]
					break
				}
			}
		}
	}

	return image.ID
}

func (image *Image) cached() (bool, error) {
	//Check with the local db
	return false, nil
}
func (image *Image) save() {
}

func analyzeLayer(endpoint, imagePath, layerID, parentLayerID, imageFormat string) error {
	payload := Layer{ID: layerID, Path: imagePath, ParentID: parentLayerID, ImageFormat: imageFormat}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", endpoint+postLayerURI, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 201 {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("Got response %d with message %s", response.StatusCode, string(body))
	}

	return nil
}

func getVulnerabilities(endpoint, layerID, minimumPriority string) ([]APIVulnerability, error) {
	response, err := http.Get(endpoint + fmt.Sprintf(getLayerVulnerabilitiesURI, layerID, minimumPriority))
	if err != nil {
		return []APIVulnerability{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return []APIVulnerability{}, fmt.Errorf("Got response %d with message %s", response.StatusCode, string(body))
	}

	var apiResponse APIVulnerabilitiesResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		return []APIVulnerability{}, err
	}

	return apiResponse.Vulnerabilities, nil
}
