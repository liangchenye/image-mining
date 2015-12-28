package imagesource

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/liangchenye/image-mining/collect"
	"github.com/liangchenye/image-mining/libs"
)

const (
	RepoCacheDir       = "/tmp/.imageMining"
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
	fmt.Println(gitFile)
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

	return images
}
