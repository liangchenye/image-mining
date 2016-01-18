package libs

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func DemoImage() Image {
	return Image{Format: "Docker", User: "docker.io", Repo: "alpine", Tag: "edge"}
}

func TestImagePull(t *testing.T) {
	image := DemoImage()
	image.Pull()
}

func TestImageID(t *testing.T) {
	image := DemoImage()

	id := image.GetID()
	fmt.Println("Get ID ", id)
}

func TestImageSave(t *testing.T) {
	image := DemoImage()

	path, err := image.Save()
	fmt.Println("Save ", path, err)
}

func TestImageHistory(t *testing.T) {
	image := DemoImage()

	layers, err := image.History()
	fmt.Println("History ", layers, err)
}

func TestImageScan(t *testing.T) {
	image := DemoImage()

	if err := image.Scan(); err != nil {
		fmt.Println(err)
	}

	vulnerabilities, err := image.GetVuln()
	if err != nil && len(vulnerabilities) == 0 {
		fmt.Println("Bravo, your image looks SAFE !")
	}
	for _, vulnerability := range vulnerabilities {
		fmt.Printf("- # %s\n", vulnerability.ID)
		fmt.Printf("  - Priority:    %s\n", vulnerability.Priority)
		fmt.Printf("  - Link:        %s\n", vulnerability.Link)
		fmt.Printf("  - Description: %s\n", vulnerability.Description)
	}
}

func TestCompress(t *testing.T) {
	uri := "test.tar"
	f, _ := os.Create(uri)
	f.Close()
	newUri, err := compressLayer(uri)
	assert.Nil(t, err)
	assert.Equal(t, newUri, uri+".gz")
	os.Remove(newUri)
}
