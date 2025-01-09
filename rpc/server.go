package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func main() {
	const address = "localhost:9001"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting tcp server")
		return
	}
	defer listener.Close()
	fmt.Printf("Listening on %s\n", address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var size uint32
	err := binary.Read(conn, binary.BigEndian, &size)
	if err != nil {
		fmt.Println("Error reading packet size")
		return
	}
	data := make([]byte, size)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		fmt.Println("Error reading packet data")
		return
	}
	original, err := deserialize(data)
	if err != nil {
		fmt.Println("Could not deserialize data: %w", err)
		return
	}
	callRPC(original)
}
