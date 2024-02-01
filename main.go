package main

import (
	"fmt"
	"sync"
)

type FileWriteMsg struct {
	Data []byte
}

func main() {
	var wg sync.WaitGroup

	//fmt.Println("main, starting server")
	//wg.Add(1)
	//go serve(5555, "tcp", &wg)

	wg.Add(1)
	go client(5555, &wg)
	wg.Wait()
	fmt.Println("fin main")

}
