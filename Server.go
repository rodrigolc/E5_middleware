package main

import (
	"./MyRPC"
)

type Echoer interface {
	Echo(line string) string
	ReverseEcho(line string) string
}

type Echo struct {
}

func (e *Echo) Echo(line string) string {
	return line
}
func (e *Echo) ReverseEcho(line string) string {

	temp := []rune(line)
	for i, j := 0, len(temp)-1; i < j; i, j = i+1, j-1 {
		temp[i], temp[j] = temp[j], temp[i]
	}
	return string(temp)
}

type EchoProxy struct {
}

func Lookup() {
	lookupAddress := "localhost:4444"
	var lookup MyRPC.LookUp = MyRPC.LookUp{}
	lookup.Init(lookupAddress)
}

func Server() {
	echo := Echo{}

	echoAddress := "localhost:5555"
	lookupAddress := "localhost:4444"

	var lookup MyRPC.LookUp = MyRPC.LookUp{}
	aor := lookup.CreateReference(echoAddress, 1)

	aor, err := lookup.Register("Echo", aor)
}
