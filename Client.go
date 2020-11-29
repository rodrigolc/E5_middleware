package main

import (
	"./MyRPC"
)

//Client Proxy
type EchoerProxy struct {
	aor       MyRPC.AbsoluteObjectReference
	requestor MyRPC.Requestor
}

func (e *EchoerProxy) Echo(line string) string {

	call := MyRPC.Call{"Echo", []interface{}{line}}
	newInv := MyRPC.Invocation{e.aor, call}
	newLine, err := e.requestor.Request(newInv)
	if err != nil {
		panic(err)
	}
	return newLine[0].(string)
}

func (e *EchoerProxy) ReverseEcho(line string) string {

	call := MyRPC.Call{"ReverseEcho", []interface{}{line}}
	newInv := MyRPC.Invocation{e.aor, call}
	newLine, err := e.requestor.Request(newInv)
	if err != nil {
		panic(err)
	}
	return newLine[0].(string)
}

func Client() {
	lookup := MyRPC.LookUpProxy{"localhost:5555"}
	var echoProxy EchoProxy
	echoProxy, err := lookup.LookUp("Echo")

}
