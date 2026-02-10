package main

import (
    "fmt"
    "time"
    "project/Network"
	"project/WorldView"
)


func main() {
    port := 40000

    tx := make(chan WorldView.WorldView)
    rx := make(chan WorldView.WorldView)

    // Start receiver goroutine
    go func() {
        err := Network.Reciever(rx, port)
        if err != nil {
            fmt.Println("Receiver error:", err)
        }
    }()

    // Start broadcaster goroutine
    go func() {
        err := Network.WordlView_broadcast(tx, port)
        if err != nil {
            fmt.Println("Broadcast error:", err)
        }
    }()

    // Simulate sending messages
    go func() {
        epoch := uint64(1)
        for {
            msg := WorldView.WorldView{
                SenderID: "TestNode",
                Epoch:    epoch,
                Msg:      fmt.Sprintf("Hello %d", epoch),
            }
            tx <- msg
            epoch++
            time.Sleep(5000 * time.Millisecond) // send every 200ms
        }
    }()

    // Print received messages
    for view := range rx {
        fmt.Printf("Received message from %s: Epoch=%d, Msg=%s\n",
            view.SenderID, view.Epoch, view.Msg)
    }
}
