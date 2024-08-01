package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

func produce(q *queue, num_items int) error {
	for i := 0; i < num_items; i++ {
		filename := fmt.Sprintf("test_images\\test%d.jpg", i)
		image_data, err := imageFileToBytes(filename)
		if err != nil {
			continue
		}

		// Create and serialize request
		var req = request{
			image_data,
			uint32(rand.Intn(400)),
			uint32(rand.Intn(400)),
			filename,
		}
		buffer := new(bytes.Buffer)
		// Because the image (and the filename) is a slice of bytes, we need to serialize them manually.
		// For the slices we do the length first, then the data.
		_ = binary.Write(buffer, binary.LittleEndian, uint64(len(req.image)))
		_ = binary.Write(buffer, binary.LittleEndian, req.image)
		_ = binary.Write(buffer, binary.LittleEndian, req.resize_width)
		_ = binary.Write(buffer, binary.LittleEndian, req.resize_height)
		_ = binary.Write(buffer, binary.LittleEndian, uint64(len(req.filename)))
		_ = binary.Write(buffer, binary.LittleEndian, []byte(req.filename))

		// Create message and add it to queue
		msg := new(message)
		msg.id = rand.Int()
		msg.payload = buffer.Bytes()
		msg.timestamp = time.Now()
		msg.priority = 1
		q.enqueue(msg)
		fmt.Printf("Queued item %d\n", i)
	}
	return nil
}
