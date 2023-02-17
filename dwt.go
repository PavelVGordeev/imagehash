package imagehash

import (
	"sort"
)

const (
	coeff1 = 0.5
	coeff2 = -0.5
)

func DWT1d(data []float64) {
	temp := make([]float64, len(data))
	half := len(data) / 2
	for i := 0; i < half; i++ {
		k := i * 2
		temp[i] = coeff1*data[k] + coeff1*data[k+1]
		temp[i+half] = coeff1*data[k] + coeff2*data[k+1]
	}
	copy(data, temp)
}

func DWT2d(data [][]float64, level int) {
	dims := len(data)
	for k := 0; k < level; k++ {
		curlvl := 1 << k
		curdims := dims / curlvl
		row := make([]float64, curdims)
		for i := 0; i < curdims; i++ {
			copy(row, data[i])
			DWT1d(row)
			copy(data[i], row)
		}
		col := make([]float64, curdims)
		for j := 0; j < curdims; j++ {
			for i := 0; i < curdims; i++ {
				col[i] = data[i][j]
			}
			DWT1d(col)
			for i := 0; i < curdims; i++ {
				data[i][j] = col[i]
			}
		}
	}
}

func iDWT1d(data []float64) {
	temp := make([]float64, len(data))
	half := len(data) / 2
	for i := 0; i < half; i++ {
		k := i * 2
		temp[k] = (coeff1*data[i] + coeff1*data[i+half]) / coeff1
		temp[k+1] = (coeff1*data[i] + coeff2*data[i+half]) / coeff1
	}
	copy(data, temp)
}

func IDWT2d(data [][]float64, level int) {
	dims := len(data)
	for k := level - 1; k >= 0; k-- {
		curlvl := 1 << k
		curdims := dims / curlvl
		col := make([]float64, curdims)
		for j := 0; j < curdims; j++ {
			for i := 0; i < curdims; i++ {
				col[i] = data[i][j]
			}
			iDWT1d(col)
			for i := 0; i < curdims; i++ {
				data[i][j] = col[i]
			}
		}
		row := make([]float64, curdims)
		for i := 0; i < curdims; i++ {
			copy(row, data[i])
			iDWT1d(row)
			copy(data[i], row)
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
	sort.Float64s(flat)
	if len(flat)%2 == 1 {
		return flat[len(flat)/2]
	} else {
		return 0.5 * (flat[len(flat)/2-1] + flat[len(flat)/2])
	}
}

func getexcerpt(data [][]float64, width uint) [][]float64 {
	excerpt := make([][]float64, width)
	for i := 0; i < int(width); i++ {
		excerpt[i] = make([]float64, width)
		copy(excerpt[i], data[i][:width])
	}
	return excerpt
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
