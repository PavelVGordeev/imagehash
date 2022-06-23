package main

import (
	"encoding/hex"
	"math/bits"
	"math/rand"
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

func (i *Imagehash) Whash(data [][]float64, hashsize int) error {
	hashsize = int(floorp2(hashsize))
	level := bits.Len(uint(len(data)/hashsize)) - 1
	DWT2d(data, level)
	eraselevel(data, level)
	IDWT2d(data, level)
	DWT2d(data, level)
	excerpt := getexcerpt(data, hashsize)
	med := median(excerpt)
	i.hash = make([]byte, hashsize)
	ctr := 0
	offset := 0
	var acc byte
	for k := 0; k < hashsize; k++ {
		for j := 0; j < hashsize; j++ {
			if excerpt[k][j] > med {
				acc ^= 1
				acc <<= 1
			}
			ctr++
			if ctr%8 == 0 && ctr != 0 {
				i.hash[offset] = acc
				offset++
				acc = 0
			}
		}
	}
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
	i := Imagehash{}
	i.Whash(data, 16)
}
