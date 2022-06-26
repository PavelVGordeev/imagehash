package main

import (
	"reflect"
	"testing"
)

func TestImagehash_FromString(t *testing.T) {
	type fields struct {
		hash []byte
	}
	tests := []struct {
		name    string
		fields  fields
		hashstr string
		wantErr bool
	}{
		{
			name:    "error hash string",
			fields:  fields{},
			hashstr: "abcdefg",
			wantErr: true,
		},
		{
			name:    "11223344 hash",
			fields:  fields{[]byte{0x11, 0x22, 0x33, 0x44}},
			hashstr: "11223344",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i Imagehash
			j := Imagehash{hash: tt.fields.hash}
			if err := i.FromString(tt.hashstr); (err != nil) != tt.wantErr {
				t.Fatalf("FromString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(i.hash, j.hash) {
				t.Errorf("FromString() error, want %v, got %v", j.hash, i.hash)
			}
		})
	}
}

func TestImagehash_String(t *testing.T) {
	type fields struct {
		hash []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Imagehash{
				hash: tt.fields.hash,
			}
			if got := i.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImagehash_Whash(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		level    int
		whash    string
		distance int
	}{
		{
			name:     "lenna.png_16×16 bits",
			file:     "lenna.png",
			level:    16,
			distance: 5,
			whash:    "cfbccfbc43f847e947fb5e7348e341e7414741c741cf40cf40ca40fe40f441f0",
		},
		{
			name:     "lenna.png_8×8 bits",
			file:     "lenna.png",
			level:    8,
			distance: 3,
			whash:    "be98bd890b0b8f8c",
		},
		{
			name:     "rust.png_16×16 bits",
			file:     "rust.png",
			level:    16,
			distance: 5,
			whash:    "fe1ffe07f607c603800180018001800180038031c0ff81ffe1ffe1ffffffffff",
		},
		{
			name:     "rust.png_8×8 bits",
			file:     "rust.png",
			level:    8,
			distance: 3,
			whash:    "f3b10101818f8fff",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Imagehash{}
			j := &Imagehash{}
			_ = j.FromString(tt.whash)
			err := i.Whash(tt.file, tt.level)
			if err != nil {
				t.Fatal("Unexpected error:", err)
			}
			d, _ := i.Distance(*j)
			if d > tt.distance {
				t.Errorf("Whash() = %v, wantErr %v distance %v", i.String(), tt.whash, d)
			}
		})
	}
}
