package fs

import (
	"image"
	"image/gif"
	"image/png"
	"os"
)

func LoadGIF(filename string) (*gif.GIF, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return gif.DecodeAll(file)
}

func LoadPNG(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(file)
}
