package main

import "fmt"

func main() {
	requests := make(chan int)

	go func() { requests <- 5 }()

	gabu := <-requests
	fmt.Println(gabu)
}
