package main

import (
	"./MyRPC"
)

func Client() {
	lookup, err := MyRPC.GetLookUp("localhost:5555")
	if err != nil {
		panic(err)
	}

	var echoer Echoer = lookup.LookUp("Echo")
	if err != nil {
		panic(1)
	}

}
