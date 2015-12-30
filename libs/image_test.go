package libs

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func DemoImage() Image {
	return Image{Format: "Docker", User: "", Repo: "erlang", Tag: "18.1"}
}

func TestImagePull(t *testing.T) {
	//image := DemoImage()
	// image.Pull()
}

func TestImageID(t *testing.T) {
	image := DemoImage()

	id := image.GetID()
	fmt.Println(id)
}

func TestImageSave(t *testing.T) {
	image := DemoImage()

	path, err := image.Save()
	fmt.Println(path, err)
}

func TestImageHistory(t *testing.T) {
	image := DemoImage()

	layers, err := image.History()
	fmt.Println(layers, err)
}

func TestImageScan(t *testing.T) {
	image := DemoImage()

	err := image.Scan()
	fmt.Println(err)

	vulnerabilities, err := image.GetVuln()
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

func TestCompress(t *testing.T) {
	uri := "test.tar"
	f, _ := os.Create(uri)
	f.Close()
	newUri, err := compressLayer(uri)
	assert.Nil(t, err)
	assert.Equal(t, newUri, uri+".gz")
	os.Remove(newUri)
}
