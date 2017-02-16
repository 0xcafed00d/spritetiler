package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
)

type ImageInfo struct {
	filename string
	img      image.Image
	size     image.Point
}

func ExitOnError(err error) {
	if err != nil {
		os.Exit(1)
	}
}

func LoadImages(filenames []string) ([]ImageInfo, error) {
	images := []ImageInfo{}

	for _, f := range os.Args[1:] {
		fmt.Print("Loading: ", f, "....")

		file, err := os.Open(f)
		if err != nil {
			return images, err
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return images, err
		}

		info := ImageInfo{f, img, img.Bounds().Size()}
		images = append(images, info)

		fmt.Println(".... Done.", info.size)
	}

	return images, nil
}

func MaxImageSize(images []ImageInfo) image.Point {
	max := image.Point{}

	for i := range images {
		if images[i].size.X > max.X {
			max.X = images[i].size.X
		}
		if images[i].size.Y > max.Y {
			max.Y = images[i].size.Y
		}
	}

	return max
}

func main() {

	images, err := LoadImages(os.Args[1:])
	ExitOnError(err)

	fmt.Println("max image size = ", MaxImageSize(images))

}
