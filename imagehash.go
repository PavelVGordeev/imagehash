package main

import (
	"encoding/hex"
	"math/rand"
)

type Imagehash struct {
	hash []byte
}

func NewHash(filename string, size uint) Imagehash {
	hash := Imagehash{hash: make([]byte, floorp2(size))}
	//hash.Whash(filename)
	return hash
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

func (i *Imagehash) Whash(data [][]float64, level int) error {
	DWT2d(data, level)
	EraseLevel(data)
	IDWT2d(data, level)
	DWT2d(data, level)
	return nil
}
func main() {
	data := make([][]float64, 1024)
	for i := 0; i < 1024; i++ {
		data[i] = make([]float64, 1024)
		for j := 0; j < 1024; j++ {
			data[i][j] = rand.Float64()
		}
	}
	Dwt2dC(data, 2)
}
