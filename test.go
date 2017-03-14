package main

import (
	"fmt"
	"github.com/lazywei/go-opencv/opencv"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

func main() {
	fImg1, _ := os.Open("IMAG0114.jpg")
	defer fImg1.Close()
	img1, _, _ := image.Decode(fImg1)

	fImg2, _ := os.Open("download.jpg")
	defer fImg2.Close()
	img2, _, _ := image.Decode(fImg2)

	rect2 := image.Rect(690, 100, 1230, 570)
	m := image.NewRGBA(img1.Bounds())
	rect := rect2.Sub(rect2.Min).Add(image.Point{500, 500})
	draw.Draw(m, m.Bounds(), img1, image.Point{0, 0}, draw.Src)
	draw.Draw(m, rect, img2, rect2.Min, draw.Src)
	toimg, _ := os.Create("new.jpg")
	defer toimg.Close()
	jpeg.Encode(toimg, m, &jpeg.Options{jpeg.DefaultQuality})

	baseImage, _, _ := image.Decode(fimg1)

	finder := facefinder.NewFinder(haarCascadeFilepath)
	faces := finder.Detect(baseImage)
	fmt.Println("asdf")
	for _, face := range faces {
		fmt.Println("Face!")
	}
}
