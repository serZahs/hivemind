package main

import (
	"fmt"
)

func main() {
	server := rpcHost{"localhost:9001"}
	procedures, err := parseProcedureFile("test.procedures")
	if err != nil {
		fmt.Println("Could not parse procedures file: %w", err)
		return
	}
	MessageBoxA := procedures[0]
	err = MessageBoxA.fill(0, "It's the final showdown", "Nuts!", 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	server.send(MessageBoxA)
}
