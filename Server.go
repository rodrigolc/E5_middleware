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

//Client Proxy
type EchoProxy struct {
	AOR       MyRPC.AbsoluteObjectReference
	Requestor MyRPC.Requestor
}

func (e *EchoProxy) Echo(line string) string {
	call := MyRPC.Call{"Echo", []interface{}{line}}
	newInv := MyRPC.Invocation{e.AOR, call}
	newLine, err := e.Requestor.Request(newInv)
	if err != nil {
		panic(err)
	}
	return newLine[0].(string)
}

func (e *EchoProxy) ReverseEcho(line string) string {
	call := MyRPC.Call{"ReverseEcho", []interface{}{line}}
	newInv := MyRPC.Invocation{e.AOR, call}
	newLine, err := e.Requestor.Request(newInv)
	if err != nil {
		panic(err)
	}
	return newLine[0].(string)
}

func Server() {
	echoAddress := "localhost:5555"
	lookupAddress := "localhost:4444"
	echo := Echo{}
	var lookup MyRPC.LookUp = MyRPC.LookUp{}
	go lookup.Init(lookupAddress)
	aor := lookup.CreateReference(echoAddress, 1)

	aor, err := lookup.Register("Echo", aor)
}
