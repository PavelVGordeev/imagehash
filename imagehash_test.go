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
	type fields struct {
		hash []byte
	}
	type args struct {
		data  [][]float64
		level int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Imagehash{
				hash: tt.fields.hash,
			}
			if err := i.Whash(tt.args.data, tt.args.level); (err != nil) != tt.wantErr {
				t.Errorf("Whash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewHash(t *testing.T) {
	type args struct {
		filename string
		size     uint
	}
	tests := []struct {
		name string
		args args
		want Imagehash
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHash(tt.args.filename, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
