package MyRPC

import (
	"fmt"
	"io/ioutil"
	"net"
)

//Server Request Handler
type ServerRequestHandler interface {
	SetUp(inv *Invoker, address string) (ServerRequestHandler, error)
	Close()
	TearDown()
	Listen() //loop principal de receber requisições
	Handle(conn net.Conn)
	Send(conn net.Conn, message []byte) (int, error)
	Receive(conn net.Conn) ([]byte, error)
}

type ServerRequestHandlerTCP struct {
	listener *net.TCPListener
	invoker  *Invoker
	close    bool
}

func (srh ServerRequestHandlerTCP) SetUp(inv *Invoker, address string) (ServerRequestHandler, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	listener, err2 := net.ListenTCP("tcp", addr)
	if err2 != nil {
		return nil, err2
	}
	srh.invoker = inv
	srh.listener = listener
	srh.close = false
	return srh, nil
}

func (srh ServerRequestHandlerTCP) TearDown() {
	srh.listener.Close()
}

func (srh ServerRequestHandlerTCP) Listen() {
	for !srh.close {
		fmt.Println("listen1")

		conn, err := srh.listener.AcceptTCP()
		fmt.Println("listen2")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Println("conn?", conn)

		go srh.Handle(conn)
	}
}
func (srh ServerRequestHandlerTCP) Handle(conn net.Conn) {
	fmt.Println("handle1", conn)

	message, err := ioutil.ReadAll(conn)
	fmt.Println("handle2", conn)

	if err != nil {
		fmt.Println(err)
	}
	message, err = (*srh.invoker).Invoke(message)
	if err != nil {
		fmt.Println(err)
	}

	_, err = srh.Send(conn, message)
	if err != nil {
		fmt.Println(err)
	}
	conn.Close()
}

func (srh ServerRequestHandlerTCP) Send(conn net.Conn, message []byte) (int, error) {
	ret, err := conn.Write(message)
	conn.Close()
	return ret, err
}

func (srh ServerRequestHandlerTCP) Receive(conn net.Conn) ([]byte, error) {
	return ioutil.ReadAll(conn)
}

func (srh ServerRequestHandlerTCP) Close() {
	srh.close = true
}

//Absolute Object Reference (AOR)
//Client Proxy
//Requestor
//Marshaller
//Invoker
//Lookup
