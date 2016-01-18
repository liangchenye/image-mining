package main

import (
	"testing"

	"github.com/liangchenye/image-mining/collect"
	_ "github.com/liangchenye/image-mining/collect/imagesource"
)

func BenchmarkImages(b *testing.B) {
	collect.LoadRepos()
	images := collect.ListImages()
	ScanImages(images)
	for i := 0; i < b.N; i++ {
		QueryImages(images)
	}
}
