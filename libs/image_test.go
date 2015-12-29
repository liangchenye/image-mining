package libs

import (
	"fmt"
	"testing"
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
}
