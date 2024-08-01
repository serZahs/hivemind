package main

import (
	"fmt"
	"sync"
)

type request struct {
	image         []byte
	resize_width  uint32
	resize_height uint32
	filename      string // Original image file
}

func main() {
	var jobs queue
	jobs.init()
	err := produce(&jobs, 30)
	if err != nil {
		fmt.Printf("Failed to produce!")
		return
	}
	fmt.Println("Done producing!")
	var wait_group sync.WaitGroup
	num_consumers := 3
	for i := 0; i < num_consumers; i++ {
		wait_group.Add(1)
		go func() {
			err = consume(&wait_group, &jobs)
			if err != nil {
				fmt.Printf("Failed to consume!")
				return
			}
		}()
	}
	wait_group.Wait()
	fmt.Println("Done consuming!")
}
