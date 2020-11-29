package main

import (
	"fmt"

	"./MyRPC"
)

func Client() {
	//echoAddress := "127.0.0.1:5555" //deve achar pelo lookup
	lookupAddress := "127.0.0.1:4146"
	fmt.Println("Oi")
	lookupProxy := MyRPC.LookUpProxy{}
	lookupProxy.New(lookupAddress)
	fmt.Println("Oi")
	echoAOR, err := lookupProxy.LookUp("Echo")
	if err != nil {
		fmt.Println("errooo")
	}
	echo := EchoProxy{}
	echo.New(echoAOR)
	fmt.Println("Oi", echo.ReverseEcho("Oi"))

}
