package willFace

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/lazywei/go-opencv/opencv"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
)

//file directories containing images of faces, and the OpenCv haarcascade
const faceDir = "Faces"
const altFaceDir = "sideFaces"
const haarCascadeDir = "/home/ubuntu/workspace/haarcascade_frontalface_alt.xml"

func DrawFace(orgImage image.Image) image.Image { //we take and return a golang image
	imageFace := orgImage

	cvImage := opencv.FromImage(imageFace)
	cascade := opencv.LoadHaarClassifierCascade(haarCascadeDir)
	faces := cascade.DetectObjects(cvImage) //where the real magic happens

	var output []image.Rectangle
	for _, value := range faces {
		output = append(output, image.Rectangle{
			image.Point{value.X(), value.Y()},
			image.Point{(value.X() + value.Width()), (value.Y() + value.Height())}})
	}
	//need blank canvas to draw new image
	canvas := image.NewRGBA(imageFace.Bounds())
	draw.Draw(canvas, imageFace.Bounds(), imageFace, imageFace.Bounds().Min, draw.Src)
	for i := range output {
		output[i] = rectResize(50.0, output[i]) //increase from face detecter to head detect...not perfect
		img1 := getFace(faceDir)
		/*		black := image.NewUniform(color.Black)
				draw.Draw(canvas, output[i], black, black.Bounds().Min, draw.Over)*/
		sizedFace := imaging.Fit(img1, output[i].Dx(), output[i].Dy(), imaging.Lanczos) //makes sure all of img1 is within the target

		if sizedFace.Bounds().Dy() < output[i].Dy() || sizedFace.Bounds().Dx() < output[i].Dx() { //need to center image
			heightDiff := sizedFace.Bounds().Dy() - output[i].Dy()
			widthDiff := sizedFace.Bounds().Dx() - output[i].Dx()
			if heightDiff < 0 && widthDiff < 0 {
				//means we need to grow our sizedFace to fit
				if heightDiff > widthDiff {
					sizedFace = imaging.Resize(img1, 0, output[i].Dy(), imaging.Lanczos)
				} else if heightDiff < widthDiff {
					sizedFace = imaging.Resize(img1, output[i].Dx(), 0, imaging.Lanczos)
				} else { //height and width difference must be the same
					sizedFace = imaging.Resize(img1, output[i].Dx(), output[i].Dy(), imaging.Lanczos)
				}
				//reset height and width
				heightDiff = sizedFace.Bounds().Dy() - output[i].Dy()
				widthDiff = sizedFace.Bounds().Dx() - output[i].Dx()
			}
			//now to center sizedFace...by trimming the target
			output[i].Min.X -= (widthDiff / 2)
			output[i].Min.Y -= (heightDiff / 2)
		}

		draw.Draw(canvas, output[i], sizedFace, imageFace.Bounds().Min, draw.Over)
	}
	if len(output) == 0 { //makes the face we are drawing 1/3 screen height
		sizedFace := imaging.Resize(getFace(altFaceDir), 0, imageFace.Bounds().Dy()/3, imaging.Lanczos)
		draw.Draw(
			canvas,
			image.Rectangle{ //this places the face in the middle of the image on the right side with half of the face offscreen
				image.Pt(canvas.Bounds().Max.X-(sizedFace.Bounds().Dx()/2), canvas.Bounds().Max.Y-2*sizedFace.Bounds().Dy()),
				canvas.Bounds().Max},
			sizedFace,
			sizedFace.Bounds().Min,
			draw.Over)
	}
	fmt.Println("New image processed faces found:", len(output))
	return canvas
}
func rectResize(sizeInc float64, rect image.Rectangle) image.Rectangle { //size Inc is percentage increase
	adj_width := int(sizeInc * float64(rect.Dx()) / 200)//half the amount of the pixels to be adjusted
	adj_height := int(sizeInc * float64(rect.Dy()) / 200)

	return image.Rect(
		rect.Min.X-adj_width,
		rect.Min.Y-adj_height,
		rect.Max.X+adj_width,
		rect.Max.Y+adj_height)
}
func getFace(dir string) image.Image {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	i := rand.Intn(len(files))
	faceFile := files[i].Name()
	fImg1, err := os.Open(path.Join(dir, faceFile))
	if err != nil {
		fmt.Printf("%s", err)
	}
	img1, _, _ := image.Decode(fImg1)
	fImg1.Close()
	if dir == faceDir {
		i = rand.Intn(2)
	} else {
		i = 1
	}
	if i == 1 {
		img1 = imaging.FlipH(img1)
	}
	return img1
}
