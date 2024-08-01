package main

import (
    "fmt"
    "net"
    "encoding/gob"
    "sync"
    "example/project/core"
)

func DeserializeTask(decoder *gob.Decoder) (core.Task, bool) {
    var task core.Task
    err := decoder.Decode(&task)
    if err != nil {
        fmt.Println(err)
        return task, false
    }
    return task, true
}

func InitWorkers(task_channel chan core.Task, quit_channel chan bool, 
    num_workers int, searches *[]int) {
    for i := 0; i < num_workers; i++ {
        go func(index int) {
            var mutex sync.Mutex
            for {
                select {
                case task := <-task_channel:
                    fmt.Printf("Worker %d got task\n", index)
                    result := core.FindBytesInArray(task.Text, task.Pattern)
                    //fmt.Println(result)

                    mutex.Lock()
                    for _, v := range result {
                        *searches = append(*searches, task.Start + v)
                    }
                    mutex.Unlock()
                case <-quit_channel:
                    //fmt.Println("Worker quitting")
                    return
                }
            }
        }(i)
    }
}

func ServeClient(connection net.Conn) {
    decoder := gob.NewDecoder(connection)
    encoder := gob.NewEncoder(connection)
    
    task_channel := make(chan core.Task)
    quit_channel := make(chan bool)
    var searches []int
    const num_workers = 8
    InitWorkers(task_channel, quit_channel, num_workers, &searches)

    for {
        //fmt.Println("Awaiting tasks...")
        task, success := DeserializeTask(decoder)
        if !success { return }

        if len(task.Text) == 0 {
            //fmt.Println("Finished")
            for i := 0; i < num_workers; i++ {
                quit_channel <- true
            }
            encoder.Encode(searches)
            connection.Close()
            return
        } else {
            task_channel <- task
        }
    }
}

func InitNode(port_string string) {
    port := ":" + port_string
    listener, err := net.Listen("tcp", port)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer listener.Close()

    for {
        connection, err := listener.Accept()
        if err != nil {
            fmt.Println(err)
            return
        }
        go ServeClient(connection)
    }
    
}

func main() {
    const port = "3000"
    InitNode(port)
}