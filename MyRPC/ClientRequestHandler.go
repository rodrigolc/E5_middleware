package MyRPC

import (
	"io/ioutil"
	"net"
)

//Client Request Handler (CRH)
type ClientRequestHandler interface {
	SendReceive(message []byte) ([]byte, error)
	SetUp(addr string) (ClientRequestHandler, error)
}

type ClientRequestHandlerTCP struct {
	ServerAddr *net.TCPAddr
}

func (crh ClientRequestHandlerTCP) SendReceive(message []byte) ([]byte, error) {
	conn, err := net.DialTCP("tcp", nil, crh.ServerAddr)
	if err != nil {
		return nil, err
	}
	conn.Write(message)
	err = conn.CloseWrite()
	if err != nil {
		return nil, err
	}
	var response []byte
	response, err = ioutil.ReadAll(conn)
	conn.Close()
	return response, err
}
func (crh ClientRequestHandlerTCP) SetUp(addr string) (*ClientRequestHandlerTCP, error) {
	serverAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	crh.ServerAddr = serverAddr
	return &crh, err
}
