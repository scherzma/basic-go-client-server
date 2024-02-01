package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func client(port int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer fmt.Println("fin defer client")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	wg.Add(1)
	time.Sleep(500 * time.Millisecond)
	go connect(port, wg, "one")

	for i := 0; i < 100; i++ {
		wg.Add(1)
		time.Sleep(1 * time.Millisecond)
		go connect(port, wg, strconv.Itoa(i))
	}
	fmt.Println("fin client")
}

func connect(port int, wg *sync.WaitGroup, message_cont string) {
	defer wg.Done()
	//defer fmt.Println("fin connect")

	conn, err := net.Dial("tcp", os.Getenv("IP1")+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	message := "Hello, Server! " + time.Now().String() + " " + message_cont
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err.Error())
		return
	}
}
