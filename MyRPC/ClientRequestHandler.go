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
	//println("CRH! sendreceive! Dial?")
	conn, err := net.DialTCP("tcp", nil, crh.ServerAddr)
	//println("CRH! sendreceive! Dial!")
	if err != nil {
		return nil, err
	}
	//println("CRH! sendreceive! write?")
	conn.Write(message)
	//println("CRH! sendreceive! write!")
	err = conn.CloseWrite()
	//println("CRH! sendreceive! closewrite!")
	if err != nil {
		return nil, err
	}
	var response []byte
	//println("CRH! sendreceive! readAll?")
	response, err = ioutil.ReadAll(conn)
	//println("CRH! sendreceive! readAll!")
	conn.Close()
	//println("CRH! sendreceive! conn close!")
	return response, err
}
func (crh ClientRequestHandlerTCP) SetUp(addr string) (ClientRequestHandler, error) {
	serverAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	crh.ServerAddr = serverAddr
	return crh, err
}
