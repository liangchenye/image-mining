package libs

import (
	"fmt"
)

type Layer struct {
	ID          string
	ParentID    string
	URL         string
	ImageFormat string
}

//Only support 'docker' at present
type Image struct {
	Format string //Docker or ACI
	User   string
	Repo   string
	Tag    string

	LayerInfo Layer
}

func (image *Image) Pull() error {
	if cached, err := image.cached(); err != nil {
		return err
	} else if cached {
		return nil
	}
	command := fmt.Sprintf("docker pull %s:%s", image.Repo, image.Tag)
	//id := fmt.Sprintf("docker get id blabla")
	fmt.Println(command)
	//using info in the contri blabla
	return nil
}

func (image *Image) Scan() {
	// scan with Clair
	//backend scan policy
	image.save()
}

func (image *Image) Clear() {
}

func (image *Image) cached() (bool, error) {
	//Check with the local db
	return true, nil
}
func (image *Image) save() {
}
