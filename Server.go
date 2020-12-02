package main

import (
	"errors"
	"fmt"
	"time"

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

func (echo *EchoProxy) New(aor MyRPC.AbsoluteObjectReference) *EchoProxy {
	*echo = EchoProxy{aor, MyRPC.Requestor{MyRPC.ClientRequestHandlerTCP{}}} //ID fixo do Echo
	return echo
}

func (e *EchoProxy) Echo(line string) string {
	call := MyRPC.Call{"Echo", []interface{}{line}}
	newInv := MyRPC.Invocation{e.AOR, call}
	newLine, err := e.Requestor.Request(newInv)
	if err != nil {
		fmt.Println(err)
	}
	return newLine[0].(string)
}

func (e *EchoProxy) ReverseEcho(line string) string {
	call := MyRPC.Call{"ReverseEcho", []interface{}{line}}
	newInv := MyRPC.Invocation{e.AOR, call}
	newLine, err := e.Requestor.Request(newInv)
	if err != nil {
		fmt.Println(err)
	}
	return newLine[0].(string)
}

type EchoInvoker struct {
	Echo *Echo
	SRH  *MyRPC.ServerRequestHandler
}

func (inv EchoInvoker) Invoke(message []byte) ([]byte, error) {
	m := MyRPC.Marshaller{}
	op := MyRPC.Invocation{}
	err := m.Unmarshal(message, &op)
	if err != nil {
		return nil, err
	}

	switch op.Call.Method {
	case "Echo":
		ech := (*inv.Echo).Echo(op.Call.Args[0].(string))
		return m.Marshal(ech)
	case "ReverseEcho":
		rev := (*inv.Echo).ReverseEcho(op.Call.Args[0].(string))
		return m.Marshal(rev) //parece errado, mas é isso mesmo
	default:
		return nil, errors.New("Operação não reconhecida")
	}
}
func LookupServer() {
	//println("lookup!")
	lookupAddress := "127.0.0.1:4146"
	var lookup MyRPC.LookUp = MyRPC.LookUp{}
	//println("init?")
	lookup.Init(lookupAddress)
}
func Server() {
	echoAddress := "127.0.0.1:5555"
	lookupAddress := "127.0.0.1:4146"
	echo := Echo{}
	var srh MyRPC.ServerRequestHandler = MyRPC.ServerRequestHandlerTCP{}
	//println("SERVER! invoker?")
	var inv MyRPC.Invoker = EchoInvoker{&echo, &srh}
	//println("SERVER! invoker! setup?")
	srh, _ = srh.SetUp(&inv, echoAddress)

	//println("lookup?")
	go LookupServer()
	//println("lookup? wait 2")
	time.Sleep(2 * time.Second)
	//println("lookup? waited 2")
	lookupProxy := MyRPC.LookUpProxy{}
	lookupProxy.New(lookupAddress)
	//println("SERVER! lookup? register?")
	lookupProxy.Register("Echo", lookupProxy.CreateReference(echoAddress, "1"))
	//println("SERVER! lookup! register!")
	//println("SERVER! setup! listen?")
	srh.Listen()
}
