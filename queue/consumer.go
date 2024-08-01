package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"
)

func consume(wait_group *sync.WaitGroup, q *queue) error {
	defer wait_group.Done()
	for {
		msg := q.dequeue()
		if msg == nil {
			break
		}
		reader := bytes.NewReader(msg.payload)

		// Deserialize the request
		var req request
		var image_size uint64
		var filename_size uint64
		_ = binary.Read(reader, binary.LittleEndian, &image_size)
		req.image = make([]byte, image_size)
		_ = binary.Read(reader, binary.LittleEndian, &req.image)
		_ = binary.Read(reader, binary.LittleEndian, &req.resize_width)
		_ = binary.Read(reader, binary.LittleEndian, &req.resize_height)
		_ = binary.Read(reader, binary.LittleEndian, &filename_size)
		filename_bytes := make([]byte, filename_size)
		_ = binary.Read(reader, binary.LittleEndian, &filename_bytes)
		req.filename = string(filename_bytes)

		// Load the image and resize it.
		img, err := imageFromBytes(req.image)
		if err != nil {
			return err
		}
		resizeImage(img, uint(req.resize_width), uint(req.resize_height), req.filename)
		fmt.Printf("Resized %s (%dx%d)\n", req.filename, req.resize_width, req.resize_height)
	}
	return nil
}
