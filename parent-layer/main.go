package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/liangchenye/image-mining/libs"
)

/*
This program is used to check if a 'layer' will have different 'parent' layers.
*/

type HistoryInfo struct {
	Parent string
	Image  libs.Image
}

var (
	historyStore = make(map[string]HistoryInfo)
	succ         int
)

func main() {
	succ = 0
	if images, err := LoadLocalImages(); err != nil {
		fmt.Println(err)
	} else if len(images) < 2 {
		fmt.Println("No local image found.")
	} else {
		for index := 1; index < len(images); index++ {
			strs := strings.Fields(images[index])
			if image, err := libs.ImageNew("Docker", "", strs[0], strs[1]); err == nil {
				StoreHistory(image)
			} else {
				fmt.Println(err)
			}
		}
	}
	fmt.Sprintf("Yes, no layer has different parents, verifed %d times", succ)
}

func LoadLocalImages() ([]string, error) {
	var stderr bytes.Buffer
	cmd := exec.Command("docker", "images")
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

	return layers, nil
}

//The first one is the image id (no parent)
func StoreHistory(image libs.Image) {
	history, err := image.History()
	if err != nil {
		return
	}
	for i := 0; i < len(history); i++ {
		if i > 0 {
			if parent, alreadyExists := historyStore[history[i]]; alreadyExists {
				if parent.Parent != history[i-1] {
					panic(fmt.Sprintf("!!! layer has different parents,%s in [%s/%s] and [%s/%s]", history[i], image.Repo, image.Tag, parent.Image.Repo, parent.Image.Tag))
				} else {
					fmt.Println("same parent")
					succ++
				}
			} else {
				var hi HistoryInfo
				hi.Parent = history[i-1]
				hi.Image = image
				historyStore[history[i]] = hi
			}
		}
	}
}
