package main

import (
	"golang.org/x/image/draw"
	"image"
)

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
