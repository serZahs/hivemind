package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func imageFileToBytes(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Could not open file %s\n", path)
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Printf("Could not get size of file %s\n", path)
		return nil, err
	}

	buffer := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Printf("Could not read file %s\n", path)
		return nil, err
	}
	return buffer, nil
}

func imageFromBytes(data []byte) (image.Image, error) {
	reader := bytes.NewReader(data)
	_, format, err := image.Decode(reader)
	if err != nil {
		fmt.Println("Could not decode image")
		return nil, err
	}
	reader.Seek(0, 0)
	if format == "jpeg" {
		image, err := jpeg.Decode(reader)
		if err != nil {
			fmt.Printf("Could not decode jpeg\n")
			return nil, err
		}
		return image, err
	}
	// TODO: Handle other formats
	return nil, err
}

func resizeImage(img image.Image, new_width uint, new_height uint, filename string) error {
	file, _ := os.Create(prefixFilePath(filename, "resized_"))
	new_img := resize.Resize(new_width, new_height, img, resize.Lanczos3)
	err := jpeg.Encode(file, new_img, nil)
	if err != nil {
		fmt.Println("Could not encode image")
	}
	return err
}

func prefixFilePath(path, prefix string) string {
	directory := filepath.Dir(path)
	filename := filepath.Base(path)
	extension := filepath.Ext(filename)
	name := filename[:len(filename)-len(extension)]
	new_filename := prefix + name + extension
	result := filepath.Join(directory, new_filename)
	return result
}
