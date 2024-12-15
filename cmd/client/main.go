package main

import (
	"fmt"
	"net"
	"time"
)

var errors int = 0
var sent int = 0

func sendThousandMessage(c net.Conn) {
	for i := 0; i < 70; i++ {
		buf := make([]byte, 1024)
		_, err := c.Write([]byte(buf))
		sent++
		if err != nil {
			errors++
			panic(err)
		}
	}
	fmt.Print("Thread ended")
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:10000")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 2; i++ {
		go sendThousandMessage(conn)
	}

	fmt.Printf("Errors: %d\n", errors)
	for {
		time.Sleep(1000 * time.Millisecond)
		fmt.Printf("Errors: %d\n", errors)
		fmt.Printf("Sent: %d\n", sent)
	}
}
