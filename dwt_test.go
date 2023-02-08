package main

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
)

const eps = 1e-5

func TestDWT1d(t *testing.T) {
	tests := []struct {
		name string
		data []float64
		want []float64
	}{
		{name: "simple vector",
			data: []float64{3.4, 8.4, 6.5, 9.9, 7.6, 9.3, 6.6, 5.2},
			want: []float64{5.9, 8.2, 8.45, 5.9, -2.5, -1.7, -0.85, 0.7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DWT1d(tt.data)
			for i := range tt.data {
				if math.Abs(tt.data[i]-tt.want[i]) > eps {
					t.Errorf("DWT1d() = %v, want %v", tt.data, tt.want)
					break
				}
			}
		})
	}
}

func BenchmarkDWT1d(b *testing.B) {
	data := []float64{3.4, 8.4, 6.5, 9.9, 7.6, 9.3, 6.6, 5.2}
	for i := 0; i < b.N; i++ {
		DWT1d(data)
	}
}

func TestDWT2d(t *testing.T) {
	tests := []struct {
		name  string
		data  [][]float64
		want  [][]float64
		level int
	}{
		{name: "simple matrix level1",
			data: [][]float64{
				{3.4, 8.4, 6.5, 9.9}, {7.6, 9.3, 6.6, 5.2},
				{5.1, 5.9, 2.3, 3.7}, {8.3, 0.9, 3.6, 0.2}},
			want: [][]float64{
				{7.175, 7.05, -1.675, -0.5}, {5.05, 2.45, 1.65, 0.5},
				{-1.275, 1.15, -0.825, -1.2}, {0.45, 0.55, -2.05, -1.2},
			},
			level: 1,
		},
		{name: "simple matrix level2",
			data: [][]float64{
				{3.4, 8.4, 6.5, 9.9}, {7.6, 9.3, 6.6, 5.2},
				{5.1, 5.9, 2.3, 3.7}, {8.3, 0.9, 3.6, 0.2}},
			want: [][]float64{
				{5.43125, 0.68125, -1.675, -0.5}, {1.68125, -0.61875, 1.65, 0.5},
				{-1.275, 1.15, -0.825, -1.2}, {0.45, 0.55, -2.05, -1.2},
			},
			level: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DWT2d(tt.data, tt.level)
		outer:
			for i := range tt.data {
				for j := range tt.data[0] {
					if math.Abs(tt.data[i][j]-tt.want[i][j]) > eps {
						t.Errorf("DWT2d() = %v, want %v", tt.data, tt.want)
						break outer
					}
				}
			}
		})
	}
}

func Test_iDWT1d(t *testing.T) {
	tests := []struct {
		name string
		data []float64
		want []float64
	}{
		{name: "simple vector",
			want: []float64{4, 0.5, 0.75, 0.2, 2, 0.6},
			data: []float64{2.25, 0.475, 1.3, 1.75, 0.275, 0.7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iDWT1d(tt.data)
			for i := range tt.data {
				if math.Abs(tt.data[i]-tt.want[i]) > eps {
					t.Errorf("iDWT1d() = %v, want %v", tt.data, tt.want)
					break
				}
			}
		})
	}
}

func Test_iDWT2d(t *testing.T) {
	tests := []struct {
		name  string
		data  [][]float64
		want  [][]float64
		level int
	}{
		{name: "simple matrix level1",
			want: [][]float64{
				{4, 0.5, 0.75, 0.2}, {0.8, 3, 1.2, 1.4},
				{0.7, 0.2, 0.1, 1}, {2, 2.5, 0.1, 0.9}},
			data: [][]float64{
				{2.075, 0.8875, 0.325, 0.0875}, {1.35, 0.525, 0, -0.425},
				{0.175, -0.4125, 1.425, 0.1875}, {-0.9, 0.025, 0.25, -0.025},
			},
			level: 1,
		},
		{name: "simple matrix level2",
			want: [][]float64{
				{4, 0.5, 0.75, 0.2}, {0.8, 3, 1.2, 1.4},
				{0.7, 0.2, 0.1, 1}, {2, 2.5, 0.1, 0.9}},
			data: [][]float64{
				{1.209375, 0.503125, 0.325, 0.0875}, {0.271875, 0.090625, 0, -0.425},
				{0.175, -0.4125, 1.425, 0.1875}, {-0.9, 0.025, 0.25, -0.025},
			},
			level: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			IDWT2d(tt.data, tt.level)
		outer:
			for i := range tt.data {
				for j := range tt.data[0] {
					if math.Abs(tt.data[i][j]-tt.want[i][j]) > eps {
						t.Errorf("iDWT2d() = %v, want %v", tt.data, tt.want)
						break outer
					}
				}
			}
		})
	}
}

func BenchmarkDWT2d(b *testing.B) {
	data := make([][]float64, 1024)
	for i := 0; i < 1024; i++ {
		data[i] = make([]float64, 1024)
		for j := 0; j < 1024; j++ {
			data[i][j] = rand.Float64()
		}
	}
	for i := 0; i < b.N; i++ {
		DWT2d(data, 2)
	}
}

func BenchmarkDWT2d_c(b *testing.B) {
	data := make([][]float64, 1024)
	for i := 0; i < 1024; i++ {
		data[i] = make([]float64, 1024)
		for j := 0; j < 1024; j++ {
			data[i][j] = rand.Float64()
		}
	}
	for i := 0; i < b.N; i++ {
		DWT2d(data, 2)
	}
}

func Test_floorp2(t *testing.T) {
	tests := []struct {
		name string
		val  int
		want uint
	}{
		{
			name: "1234",
			val:  1234,
			want: 1024,
		},
		{
			name: "32",
			val:  32,
			want: 32,
		},
		{
			name: "9999",
			val:  9999,
			want: 8192,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := floorp2(tt.val); got != tt.want {
				t.Errorf("floorp2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_flatten(t *testing.T) {
	tests := []struct {
		name string
		val  [][]float64
		want []float64
	}{
		{
			name: "2×2",
			val:  [][]float64{{1., 2.}, {3., 4.}},
			want: []float64{1., 2., 3., 4.},
		},
		{
			name: "3×3",
			val:  [][]float64{{1., 2., 3.}, {4., 5., 6.}, {7., 8., 9.}},
			want: []float64{1., 2., 3., 4., 5., 6., 7., 8., 9.},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := flatten(tt.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("flatten() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_median(t *testing.T) {
	tests := []struct {
		name string
		val  [][]float64
		want float64
	}{
		{
			name: "leftmost square even",
			val:  [][]float64{{8., 2., 3., 4.}, {5., 6., 7., 1.}, {9., 10., 11., 12.}, {13., 14., 15., 16.}},
			want: 8.5,
		},
		{
			name: "leftmost square odd",
			val:  [][]float64{{5., 2., 3.}, {4., 1., 6.}, {7., 8., 9.}},
			want: 5.0,
		},
		{
			name: "rightmost jagged odd",
			val:  [][]float64{{1., 2., 3., 4.}, {9.}, {6., 7., 8., 5.}},
			want: 5.,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := median(tt.val); got != tt.want {
				t.Errorf("median() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eraselevel(t *testing.T) {
	tests := []struct {
		name  string
		val   [][]float64
		want  [][]float64
		level int
	}{
		{
			name:  "level1_4×4",
			val:   [][]float64{{8., 2., 3., 4.}, {5., 6., 7., 1.}, {9., 10., 11., 12.}, {13., 14., 15., 16.}},
			want:  [][]float64{{0., 0., 3., 4.}, {0., 0., 7., 1.}, {9., 10., 11., 12.}, {13., 14., 15., 16.}},
			level: 1,
		},
		{
			name:  "level2_4×4",
			val:   [][]float64{{8., 2., 3., 4.}, {5., 6., 7., 1.}, {9., 10., 11., 12.}, {13., 14., 15., 16.}},
			want:  [][]float64{{0., 2., 3., 4.}, {5., 6., 7., 1.}, {9., 10., 11., 12.}, {13., 14., 15., 16.}},
			level: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if eraselevel(tt.val, tt.level); !reflect.DeepEqual(tt.val, tt.want) {
				t.Errorf("erase() = %v, want %v", tt.val, tt.want)
			}
		})
	}
}
