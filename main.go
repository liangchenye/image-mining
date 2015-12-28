package main

import (
	"fmt"

	"github.com/liangchenye/image-mining/collect"
	_ "github.com/liangchenye/image-mining/collect/imagesource"
)

func main() {
	collect.LoadRepos()
	fmt.Println("hello")
}
