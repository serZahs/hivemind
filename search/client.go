package main

import (
	"os"
	"fmt"
	"net"
	"time"
	"encoding/gob"
	"example/project/core"
)

func SerializeTask(encoder *gob.Encoder, task core.Task) bool {
	err := encoder.Encode(task)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func InitTCPConnection(port string) (net.Conn, bool) {
	address := "localhost:" + port
	connection, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return connection, true
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Use: client {file_name} {pattern}")
		return
	}
	filename := os.Args[1]
	pattern := os.Args[2]
	bytes, err := os.ReadFile(filename)
	if err != nil {
        fmt.Println(err)
        return
    }

    const port = "3000"

    connection, success := InitTCPConnection(port)
    if !success { return }

	encoder := gob.NewEncoder(connection)
	decoder := gob.NewDecoder(connection)

	time_start := time.Now()
	const num_streams = 8
	streams, indices := core.SplitIntoChunks(bytes, num_streams)
	for i := 0; i < num_streams; i++ {
		task := core.Task{indices[i], []byte(pattern), streams[i]}
		success = SerializeTask(encoder, task)
		if !success { return }
		//fmt.Println("Serialized task")
	}
	success = SerializeTask(encoder, core.Task{0, nil, nil}) // Add the empty task to signify the end
	if !success { return }

	var searches []int
	decoder.Decode(&searches)
    elapsed := time.Since(time_start)

    fmt.Println(searches)
    fmt.Printf("Took %v, %d results\n", elapsed, len(searches))
}