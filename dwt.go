package main

import (
	"math/bits"
	"sort"
)

// Коэффициенты вейвлетного преобразования Хаара
const (
	coeff1 = 0.5
	coeff2 = -0.5
)

// Прямое вейвлетное пробразование Хаара для вектора
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

// Прямое вейвлетное пробразование Хаара для 2D-матрицы, level - коэффициент сжатия, матрица сжимается в 2^level раз
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

// Обратное вейвлетное пробразование Хаара для вектора
func iDWT1d(data []float64) {
	temp := make([]float64, len(data))
	half := len(data) >> 1
	for i := 0; i < half; i++ {
		k := i << 1
		temp[k] = (coeff1*data[i] + coeff1*data[i+half]) / coeff1
		temp[k+1] = (coeff1*data[i] + coeff2*data[i+half]) / coeff1
	}
	copy(data, temp)
}

// Обратное вейвлетное пробразование Хаара для вектора
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

// Удаление аппроксимации, сохраняются только коэффициенты восстановления
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

// Расчет N,  <= val, являющегося степенью двойки
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

//Вспомогательная функция для расчета медианы
//func _mediansel(list []float64, idx int) float64 {
//	r := rand.Intn(len(list))
//	var (
//		low    []float64
//		high   []float64
//		pivots []float64
//	)
//	pivot := list[r]
//	for i := 0; i < len(list); i++ {
//		if list[i] < pivot {
//			low = append(low, list[i])
//		} else if list[i] == pivot {
//			pivots = append(pivots, list[i])
//		} else {
//			high = append(high, list[i])
//		}
//	}
//	if idx < len(low) {
//		return _mediansel(low, idx)
//	} else if idx < len(low)+len(pivots) {
//		return pivots[0]
//	} else {
//		return _mediansel(high, idx-len(low)-len(pivots))
//	}
//}

// Получение первого квадранта матрицы размерности width
func getexcerpt(data [][]float64, width int) [][]float64 {
	excerpt := make([][]float64, width)
	for i := 0; i < width; i++ {
		excerpt[i] = make([]float64, width)
		copy(excerpt[i], data[i][:width])
	}
	return excerpt
}
