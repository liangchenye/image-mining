package imagesource

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/liangchenye/image-mining/collect"
	"github.com/liangchenye/image-mining/libs"
)

const (
	//My data volume
	RepoCacheDir       = "/tmp/image-data/.imageMining"
	DockerOfficialRepo = "https://github.com/docker-library/official-images.git"
	DockerRepoName     = "official-images"
	DockerOfficialDir  = "library"
)

type DockerOfficialImage struct{}

func init() {
	collect.RegisterImageSource("docker-official", &DockerOfficialImage{})
}

func (doi *DockerOfficialImage) Load() bool {
	p, err := os.Stat(RepoCacheDir)
	if err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(RepoCacheDir, 0777)
		}
	} else if !p.IsDir() {
		os.Remove(RepoCacheDir)
		os.MkdirAll(RepoCacheDir, 0777)
	}

	gitFile := path.Join(RepoCacheDir, DockerRepoName, ".git")
	if _, err = os.Stat(gitFile); err != nil {
		c := exec.Command("/bin/sh", "-c", fmt.Sprintf("git clone %s", DockerOfficialRepo))
		c.Dir = RepoCacheDir
		c.Run()
	} else {
		c := exec.Command("/bin/sh", "-c", "git update")
		c.Dir = path.Join(RepoCacheDir, DockerRepoName)
		c.Run()
	}

	return true
}

func (doi *DockerOfficialImage) ListImages() (images []libs.Image) {
	officialRepoDir := path.Join(RepoCacheDir, DockerRepoName, DockerOfficialDir)
	files, _ := ioutil.ReadDir(officialRepoDir)

	for _, file := range files {
		uri := path.Join(officialRepoDir, file.Name())
		images = append(images, readImagesFromFile(uri)...)
	}

	return images
}

func readImagesFromFile(uri string) (images []libs.Image) {
	f, err := os.Open(uri)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')

		if err != nil || io.EOF == err {
			break
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		strs := strings.Split(line, ":")
		if len(strs) > 1 {
			if image, err := libs.ImageNew("Docker", "", path.Base(uri), strs[0]); err == nil {
				images = append(images, image)
			}
		}
	}
	return images
}
