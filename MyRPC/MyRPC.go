package MyRPC

import (
	"net"
)

//Client Request Handler (CRH)
type ClientRequestHandler interface {
	Send(message []byte) (int, error)
	Receive(message []byte) (int, error)
	Dial(network, addr string) (ClientRequestHandler, error)
}

type ClientRequestHandlerTCP struct {
	ServerSocket *net.TCPConn
	ServerAddr   *net.TCPAddr
}

func (crh ClientRequestHandlerTCP) Send(message []byte) (int, error) {
	return crh.ServerSocket.Write(message)
}
func (crh ClientRequestHandlerTCP) Receive(message []byte) (int, error) {
	return crh.ServerSocket.Read(message)
}
func (crh ClientRequestHandlerTCP) Dial(network, addr string) (ClientRequestHandler, error) {
	var err error
	crh.ServerAddr, err = net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}
	crh.ServerSocket, err = net.DialTCP(network, crh.ServerAddr, nil)
	if err != nil {
		return nil, err
	}
	return crh, err
}

//Server Request Handler
type ServerRequestHandler interface {
	Send(message []byte) int, error
	Receive(message []byte) int, error
	Listen(port int) net.Listener, error
	Accept() (Conn,error)
}

type ServerRequestHandlerTCP struct {
	Listener net.TCPListener

}

//Absolute Object Reference (AOR)
//Client Proxy
//Requestor
//Marshaller
//Invoker
//Lookup
