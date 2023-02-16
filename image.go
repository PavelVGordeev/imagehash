package main

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

var ErrNotImplemented = errors.New("not implemented image decoder")

type Image struct {
	path string
}

func (i Image) LoadImage() (image.Image, error) {
	file, err := os.Open(i.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var decFunc func(reader io.Reader) (image.Image, error)
	switch filepath.Ext(strings.ToLower(i.path)) {
	case ".png":
		decFunc = png.Decode
	case ".jpg", ".jpeg":
		decFunc = jpeg.Decode
	default:
		return nil, ErrNotImplemented
	}
	img, err := decFunc(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (i Image) Vectorize() ([][]float64, error) {
	img, err := i.LoadImage()
	if err != nil {
		return nil, err
	}
	return grayscale(img)
}

func grayscale(img image.Image) ([][]float64, error) {

	bounds := img.Bounds()
	imgScale := int(floorp2(min(bounds.Max.X, bounds.Max.Y)))
	resized := image.NewRGBA(image.Rect(0, 0, imgScale, imgScale))
	draw.ApproxBiLinear.Scale(resized, resized.Rect, img, bounds, draw.Over, nil)

	data := make([][]float64, imgScale)
	for y := 0; y < imgScale; y++ {
		data[y] = make([]float64, imgScale)
		for x := 0; x < imgScale; x++ {
			pixel := resized.At(x, y)
			r, g, b, _ := pixel.RGBA()
			data[y][x] = (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 65535.0
		}
	}
	return data, nil
}
