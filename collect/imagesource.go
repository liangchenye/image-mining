package collect

import (
	"fmt"
	"sync"

	"github.com/liangchenye/image-mining/libs"
)

type ImageSource interface {
	Load() bool
	ListImages() []libs.Image
}

var (
	imageSourceLock sync.Mutex
	imageSource     = make(map[string]ImageSource)
)

// RegisterImageSource provides a way to dynamically register an implementation of a
// ImageSource.
//
// If RegisterImageSource is called twice with the same name if ImageSource is nil,
// or if the name is blank, it panics.
func RegisterImageSource(name string, f ImageSource) {
	if name == "" {
		panic("Could not register a ImageSource with an empty name")
	}
	if f == nil {
		panic("Could not register a nil ImageSource")
	}

	imageSourceLock.Lock()
	defer imageSourceLock.Unlock()

	if _, alreadyExists := imageSource[name]; alreadyExists {
		panic(fmt.Sprintf("Detector '%s' is already registered", name))
	}
	imageSource[name] = f
}

func LoadRepos() {
	for _, im := range imageSource {
		im.Load()
	}
}

func ListImages() (images []libs.Image) {
	for _, im := range imageSource {
		images = append(images, im.ListImages()...)
	}

	return images
}
