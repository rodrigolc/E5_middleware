package MyRPC

import (
	"errors"
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
		conn, err := srh.listener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
		}
		go srh.Handle(conn)
	}
}
func (srh ServerRequestHandlerTCP) Handle(conn net.Conn) {
	message, err := ioutil.ReadAll(conn)
	if err != nil {
		panic(errors.New("Erro lendo no Handle"))
	}
	message, err = (*srh.invoker).Invoke(message)
	if err != nil {
		panic(errors.New("Erro no Invoke()"))
	}

	l, err := srh.Send(conn, message)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Erro no Send() - numero %d", l)))
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
