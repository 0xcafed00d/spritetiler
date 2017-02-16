package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

type ImageInfo struct {
	filename string
	img      image.Image
	size     image.Point
}

func ExitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func LoadImages(filenames []string) ([]ImageInfo, error) {
	images := []ImageInfo{}

	for _, f := range filenames {
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

func TileOnImage(images []ImageInfo, tileSz image.Point, tileCols int) image.Image {
	tileRows := len(images) / tileCols
	if len(images)%tileCols > 0 {
		tileRows++
	}

	dest := image.NewNRGBA(image.Rect(0, 0, tileSz.X*tileCols, tileSz.Y*tileRows))

	for i := range images {
		x := (i%tileCols)*tileSz.X + (tileSz.X-images[i].size.X)/2
		y := (i/tileCols)*tileSz.Y + (tileSz.Y-images[i].size.Y)/2
		draw.Draw(dest, image.Rect(x, y, tileSz.X+x, tileSz.Y+y), images[i].img, image.ZP, draw.Src)
	}

	return dest
}

type Config struct {
	Help           bool
	TileCols       int
	OutputFilename string
}

var config Config

func init() {
	flag.BoolVar(&config.Help, "h", false, "display help")
	flag.IntVar(&config.TileCols, "c", 3, "Tile column count")
	flag.StringVar(&config.OutputFilename, "o", "out.png", "output filename")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: spritetiler [options] sourcefiles")
		flag.PrintDefaults()
	}
}

func main() {

	flag.Parse()

	if len(flag.Args()) == 0 && config.Help {
		flag.Usage()
		os.Exit(1)
	}

	images, err := LoadImages(flag.Args())
	ExitOnError(err)

	maxSz := MaxImageSize(images)

	fmt.Println("max image size = ", maxSz)

	all := TileOnImage(images, maxSz, config.TileCols)

	f, err := os.Create(config.OutputFilename)
	ExitOnError(err)

	err = png.Encode(f, all)
	ExitOnError(err)
}
