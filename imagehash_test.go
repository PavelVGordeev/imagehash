package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			assert.Equal(t, i.hash, j.hash)
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
		{
			name:   "short string",
			fields: fields{[]byte{0x11, 0x22, 0x33, 0x44}},
			want:   "11223344",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Imagehash{
				hash: tt.fields.hash,
			}
			assert.Equal(t, tt.want, i.String())
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
			distance: 3,
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
			distance: 3,
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
			j := Imagehash{}
			err := j.FromString(tt.whash)
			assert.NoError(t, err)
			err = i.Whash(tt.file, tt.level)
			assert.NoError(t, err)
			d, err := i.Distance(j)
			assert.NoError(t, err)
			assert.Less(t, d, tt.distance)
		})
	}
}

func TestImagehash_Distance(t *testing.T) {
	type fields struct {
		hash []byte
	}
	type args struct {
		other Imagehash
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "equal distance=0",
			fields:  fields{[]byte{0x11, 0x22, 0x33, 0x44}},
			args:    args{other: Imagehash{[]byte{0x11, 0x22, 0x33, 0x44}}},
			want:    0,
			wantErr: assert.NoError,
		},
		{
			name:    "not equal distance=6",
			fields:  fields{[]byte{0x11, 0x22, 0x33, 0x44}},
			args:    args{other: Imagehash{[]byte{0xFF, 0x22, 0x33, 0x44}}},
			want:    6,
			wantErr: assert.NoError,
		},
		{
			name:    "incompatible hashes",
			fields:  fields{[]byte{0x11, 0x22, 0x33, 0x44}},
			args:    args{other: Imagehash{[]byte{0x11, 0x22, 0x33, 0x44, 0x55}}},
			want:    0,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Imagehash{
				hash: tt.fields.hash,
			}
			got, err := i.Distance(tt.args.other)
			if !tt.wantErr(t, err) {
				t.Fatalf("Distance() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
