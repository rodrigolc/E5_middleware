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
	//println("SRH! Setup!")
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		//println("SRH! Setup! resolve address! error")
		return nil, err
	}
	//println("SRH! Setup! Listen?", address, addr.String())
	listener, err2 := net.ListenTCP("tcp", addr)

	if err2 != nil {
		//println("SRH! Setup! Listen! error")
		return nil, err2
	}
	//println("SRH! Setup! Listen!")
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
		//println("SRH! Listen! ACCEPT???", srh.invoker, srh.listener)

		conn, err := srh.listener.AcceptTCP()

		//println("SRH! Listen! accept!", conn, err)
		if err != nil {
			//println("SRH! Listen! accept error!", conn, err)
			fmt.Println(err)
			panic(err)
		}
		fmt.Println("conn?", conn)

		go srh.Handle(conn)
	}
}
func (srh ServerRequestHandlerTCP) Handle(conn net.Conn) {
	
	//println("SRH! Handle! Message? readall")
	message, err := ioutil.ReadAll(conn)
	//println("SRH! Handle! Message!", string(message), err)

	if err != nil {
		//println("SRH! Handle! Message! error", string(message), err)
		fmt.Println(err)
	}
	//println("SRH! Handle! invoke!")
	message, err = (*srh.invoker).Invoke(message)
	if err != nil {
		//println("SRH! Handle! invoke! error")
		fmt.Println(err)
	}
	//println("SRH! Handle! send?", string(message), err)
	_, err = srh.Send(conn, message)
	//println("SRH! Handle! send!", message, err)
	if err != nil {
		//println("SRH! Handle! send! error", message, err)
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
