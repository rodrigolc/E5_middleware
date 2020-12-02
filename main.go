package main

import "time"

func main() {
	go Server()
	time.Sleep(5 * time.Second)
	Client()
}
