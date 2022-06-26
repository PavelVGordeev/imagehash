package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"math/bits"
	"os"
)

type Imagehash struct {
	hash []byte
}

func (i Imagehash) String() string {
	return hex.EncodeToString(i.hash)
}

func (i *Imagehash) FromString(hashstr string) error {
	buf, err := hex.DecodeString(hashstr)
	if err != nil {
		return err
	}
	if i != nil {
		i.hash = buf
	} else {
		i = &Imagehash{hash: buf}
	}
	return nil
}

//Поиск расстояния Хэмминга для двух хэшей идентичной длины
func (i Imagehash) Distance(other Imagehash) (int, error) {
	hamming := 0
	if len(i.hash) != len(other.hash) {
		return 0, errors.New("hashes have unequal sizes")
	}
	for idx := 0; idx < len(i.hash); idx++ {
		hamming += bits.OnesCount8(i.hash[idx] ^ other.hash[idx])
	}
	return hamming, nil
}

func (i *Imagehash) Whash(filename string, hashsize int) error {
	data, err := grayscale(filename)
	if err != nil {
		return err
	}
	hashsize = int(floorp2(hashsize))
	level := bits.Len(uint(len(data))) - 1
	hashlevel := bits.Len(uint(hashsize)) - 1
	DWT2d(data, level)
	eraselevel(data, level)
	IDWT2d(data, level)
	DWT2d(data, hashlevel)
	excerpt := getexcerpt(data, hashsize)
	med := median(excerpt)
	i.hash = make([]byte, hashsize*2)
	ctr := 0
	offset := 0
	var acc byte
	for k := 0; k < hashsize; k++ {
		for j := 0; j < hashsize; j++ {
			if excerpt[k][j] > med {
				acc ^= 1
			}
			ctr++
			if ctr%8 == 0 {
				i.hash[offset] = acc
				offset++
				acc = 0
			}
			acc <<= 1
		}
	}
	return nil
}
func main() {
	i := Imagehash{}
	i.Whash("e:\\rust.png", 16)
	fmt.Println(i)
}

func grayscale(filename string) ([][]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	imgScale := int(floorp2(min(bounds.Max.X, bounds.Max.Y)))
	res := image.NewRGBA(image.Rect(0, 0, imgScale, imgScale))
	draw.NearestNeighbor.Scale(res, res.Rect, img, bounds, draw.Over, nil)

	data := make([][]float64, imgScale)
	for y := 0; y < imgScale; y++ {
		data[y] = make([]float64, imgScale)
		for x := 0; x < imgScale; x++ {
			pixel := res.At(x, y)
			r, g, b, _ := pixel.RGBA()
			l := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			data[y][x] = l / 65536.0
		}
	}
	return data, nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
