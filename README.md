# Go Wavelet hash
Go implementation of wavelet hash using discrete wavelet transformation removing the lowest low level frequency using Haar wavelet

hashsize must be a power of 2 and less or equal minimum of image scale

Examples
===
```go
file, err := os.Open("lenna.png")
if err != nil {
    return
}
defer file.Close()
img, _ := png.Decode(file)
hash := Imagehash{}
hash.Whash(img, 8)
fmt.Println(hash.String())
hash.Whash(img, 16)
fmt.Println(hash.String())
```
output will be
```
be98bd890b0b8f8c
cfbccfbc43f847e947fb5e7348e341e7414741c741cf40cf40ca40fe40f441f0
```
Hash comparison
```go

file, err := os.Open("lenna.png")
if err != nil {
return
}
defer file.Close()
file2, err := os.Open("gopher.png")
if err != nil {
return
}
defer file2.Close()
lenna, _ := png.Decode(file)
gopher, _ := png.Decode(file2)
hash1 := Imagehash{}
hash2 := Imagehash{}
hash1.Whash(lenna, 16)
fmt.Println("lenna.png hash:", hash1.String())
hash2.Whash(gopher, 16)
fmt.Println("gopher.png hash:", hash2.String())
dist, _ := hash1.Distance(hash2)
fmt.Println("The distance between Lenna and Gopher is", dist)
```
output:
```
lenna.png hash: cfbccfbc43f847e947fb5e7348e341e7414741c741cf40cf40ca40fe40f441f0
gopher.png hash: 01800fa01ff03cf03ff83ffc1ffc1ffc0ffc07fc07fc07fc07fe07f003800200
The distance between Lenna and Gopher is 122
```
Links
___
Detailed algorithm description
https://fullstackml.com/wavelet-image-hash-in-python-3504fdd282b5