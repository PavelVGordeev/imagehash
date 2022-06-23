package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"sync"
	"time"
)

const (
	w0 = 0.5
	w1 = -0.5
	s0 = 0.5
	s1 = 0.5
)

type Data struct {
	mu   sync.RWMutex
	data [][]float64
}

func (d *Data) GetSlice(n int, axis bool) []float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	sl := make([]float64, len(d.data))
	if !axis {
		copy(sl, d.data[n])
	} else {
		for i := 0; i < len(d.data[0]); i++ {
			sl[i] = d.data[i][n]
		}
	}
	return sl

}
func (d *Data) PutSlice(data []float64, n int, axis bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !axis {
		copy(d.data[n], data)
	} else {
		for i := 0; i < len(d.data[0]); i++ {
			d.data[i][n] = data[i]
		}
	}
}
func NewData(data [][]float64) *Data {
	d := Data{
		mu:   sync.RWMutex{},
		data: make([][]float64, len(data)),
	}
	for i := 0; i < len(data); i++ {
		d.data[i] = make([]float64, len(data[i]))
		copy(d.data[i], data[i])
	}
	return &d
}

func DWT1d(data []float64) {
	temp := make([]float64, len(data))
	half := len(data) >> 1
	for i := 0; i < half; i++ {
		k := i << 1
		temp[i] = s0*data[k] + s1*data[k+1]
		temp[i+half] = w0*data[k] + w1*data[k+1]
	}
	copy(data, temp)
}

func Dwt2dC(data [][]float64, level int) {
	var (
		rows, cols int
	)
	wg := sync.WaitGroup{}
	rows = len(data)
	cols = len(data[0])
	d := NewData(data)
	for k := 0; k < level; k++ {
		curlvl := 1 << k
		curcols := cols / curlvl
		currows := rows / curlvl
		start := time.Now()
		for i := 0; i < currows; i++ {
			wg.Add(1)
			r := i
			go func() {
				row := make([]float64, curcols)
				defer wg.Done()
				copy(row, d.GetSlice(r, false))
				DWT1d(row)
				d.PutSlice(row, r, false)
			}()
		}
		wg.Wait()
		fmt.Println("time to transform concurrent rowwise", time.Since(start))
		start = time.Now()
		for j := 0; j < curcols; j++ {
			col := make([]float64, currows)
			wg.Add(1)
			c := j
			go func() {
				defer wg.Done()
				col = d.GetSlice(c, true)
				DWT1d(col)
				d.PutSlice(col, c, true)
			}()
		}
		wg.Wait()
		fmt.Println("time to transform concurrent columnwise", time.Since(start))
	}
}

func DWT2d(data [][]float64, level int) {
	var (
		rows, cols int
	)
	rows = len(data)
	cols = len(data[0])
	for k := 0; k < level; k++ {
		curlvl := 1 << k
		curcols := cols / curlvl
		currows := rows / curlvl
		row := make([]float64, curcols)
		for i := 0; i < currows; i++ {
			copy(row, data[i])
			DWT1d(row)
			copy(data[i], row)
		}
		col := make([]float64, currows)
		for j := 0; j < curcols; j++ {
			for i := 0; i < currows; i++ {
				col[i] = data[i][j]
			}
			DWT1d(col)
			for i := 0; i < currows; i++ {
				data[i][j] = col[i]
			}
		}
	}
}

func iDWT1d(data []float64) {
	temp := make([]float64, len(data))
	half := len(data) >> 1
	for i := 0; i < half; i++ {
		k := i << 1
		temp[k] = (s0*data[i] + w0*data[i+half]) / w0
		temp[k+1] = (s1*data[i] + w1*data[i+half]) / s0
	}
	copy(data, temp)
}

func IDWT2d(data [][]float64, level int) {
	var (
		rows, cols int
	)
	rows = len(data)
	cols = len(data[0])
	for k := level - 1; k >= 0; k-- {
		curlvl := 1 << k
		curcols := cols / curlvl
		currows := rows / curlvl
		col := make([]float64, currows)
		for j := 0; j < curcols; j++ {
			for i := 0; i < currows; i++ {
				col[i] = data[i][j]
			}
			iDWT1d(col)
			for i := 0; i < currows; i++ {
				data[i][j] = col[i]
			}
		}
		row := make([]float64, curcols)
		for i := 0; i < currows; i++ {
			copy(row, data[i])
			iDWT1d(row)
			copy(data[i], row)
		}
	}
}

func eraselevel(data [][]float64, level int) {
	maxlvl := bits.Len(floorp2(len(data))) - 1
	if level > maxlvl {
		return
	}
	eraseidx := 1 << (maxlvl - level)
	for i := 0; i < eraseidx; i++ {
		for j := 0; j < eraseidx; j++ {
			data[i][j] = 0.0
		}
	}
}

func floorp2(val int) uint {
	val |= val >> 1
	val |= val >> 2
	val |= val >> 4
	val |= val >> 8
	val |= val >> 16
	return uint(val - (val >> 1))
}

func flatten(data [][]float64) []float64 {
	flat := make([]float64, len(data)*len(data))
	offset := 0
	for _, row := range data {
		copy(flat[offset:offset+len(row)], row)
		offset += len(row)
	}
	return flat
}

func median(data [][]float64) float64 {
	flat := flatten(data)
	if len(flat)%2 == 1 {
		return mediansel(flat, len(flat)/2)
	} else {
		return 0.5 * (mediansel(flat, len(flat)/2-1) + mediansel(flat, len(flat)/2))
	}

}

func mediansel(list []float64, idx int) float64 {
	r := rand.Intn(len(list))
	var (
		low    []float64
		high   []float64
		pivots []float64
	)
	pivot := list[r]
	for i := 0; i < len(list); i++ {
		if list[i] < pivot {
			low = append(low, list[i])
		} else if list[i] == pivot {
			pivots = append(pivots, list[i])
		} else {
			high = append(high, list[i])
		}
	}
	if idx < len(low) {
		return mediansel(low, idx)
	} else if idx < len(low)+len(pivots) {
		return pivots[0]
	} else {
		return mediansel(high, idx-len(low)-len(pivots))
	}
}

func getexcerpt(data [][]float64, width int) [][]float64 {
	excerpt := make([][]float64, width)
	for i := 0; i < width; i++ {
		excerpt[i] = make([]float64, width)
		copy(excerpt[i], data[i][:width])
	}
	return excerpt
}
