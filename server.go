package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)

func serve(port int, proto string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("serve")
	fileWriteChannel := make(chan FileWriteMsg)

	wg.Add(1)
	go server(port, proto, fileWriteChannel, wg)

	wg.Add(1)
	go writeFile("output.txt", fileWriteChannel, wg)

}

func server(port int, proto string, fileWriteChannel chan FileWriteMsg, wg *sync.WaitGroup) {
	defer wg.Done()

	ln, err := net.Listen(proto, ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return
	}
	defer ln.Close()

	for {
		fmt.Println("Waiting for connection.")
		conn, err := ln.Accept()
		fmt.Println("Connected to: " + ln.Addr().String())
		if err != nil {
			fmt.Println("Error Accept", err.Error())
			continue
		}
		wg.Add(1)
		go handleConnection(conn, fileWriteChannel, wg)
	}
}

func handleConnection(conn net.Conn, fileWriteChannel chan FileWriteMsg, wg *sync.WaitGroup) {
	defer conn.Close()
	defer fmt.Println("fin handleConnection")
	defer wg.Done()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading", err.Error())
		return
	}

	fileWriteChannel <- FileWriteMsg{Data: buf[:n]}
}

func writeFile(filename string, fileWriteChannel chan FileWriteMsg, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error creating file", err.Error())
		return
	}
	defer file.Close()

	for msg := range fileWriteChannel {
		dataWithNewline := append(msg.Data, '\n')
		_, err := file.Write(dataWithNewline)
		if err != nil {
			fmt.Println("Error writing to file", err.Error())
			continue
		}
	}
}
