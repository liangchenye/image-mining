package main

import (
	"fmt"

	"github.com/liangchenye/image-mining/collect"
	_ "github.com/liangchenye/image-mining/collect/imagesource"
	"github.com/liangchenye/image-mining/libs"
)

func main() {
	collect.LoadRepos()
	images := collect.ListImages()
	ScanImages(images)

	QueryImages(images)
}

func ScanImages(images []libs.Image) {
	for _, image := range images {
		fmt.Println("Start to pull: ", image)
		if err := image.Pull(); err != nil {
			fmt.Println("Failed to pull: ", image)
			continue
		}
		fmt.Println("Succ in pull: ")
		fmt.Println("Start to scan: ")
		if err := image.Scan(); err != nil {
			fmt.Println("Failed to scan")
		}
		if vulns, err := image.GetVuln(); err != nil {
			fmt.Println("Failed to get vulns")
		} else {
			if len(vulns) > 0 {
				fmt.Println("Critical ", len(vulns), " found in ", image)
			}
		}
		image.Clear()
	}
}

//analyse after scanImages
func QueryImages(images []libs.Image) {
	for _, image := range images {
		vs, _ := image.GetVuln()
		fmt.Println(image.Repo, image.Tag, len(vs), vs)
	}
}
