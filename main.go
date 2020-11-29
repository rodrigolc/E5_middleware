package main

import "fmt"

func main() {
	go Server()
	var aux int
	fmt.Scan(aux)

	Client()
}
