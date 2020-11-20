package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Args struct {
	operation int
	candidate int
}

var candidates [3]int

func registerVote(candidate int) {
	candidates[candidate-1] += 1
}

func getVotes() string {
	msg := "\n\n A: " + strconv.Itoa(candidates[0]) + " B: " + strconv.Itoa(candidates[1]) + " C: " + strconv.Itoa(candidates[2]) + " - UDP\n\n;"
	return msg
}

func TranformData(requisition []byte) Args {
	var data = string(requisition)
	var dataArray = strings.Split(strings.Split(data, ";")[0], ":")
	operation, _ := strconv.Atoi(dataArray[0])
	candidate, _ := strconv.Atoi(dataArray[1])
	args := Args{operation, candidate}
	return args
}

func SRHUdp(conn *net.UDPConn, addr *net.UDPAddr, requisition []byte) {
	var args = TranformData(requisition)
	if args.operation == 1 {
		registerVote(args.candidate)
		_, err := conn.WriteToUDP([]byte("\n\n Seu Voto foi computado com sucesso UDP \n\n;"), addr)
		if err != nil {
			fmt.Print(err)
		}
	} else {
		_, err := conn.WriteToUDP([]byte(getVotes()), addr)
		if err != nil {
			fmt.Print(err)
		}
	}
}

func Server() {
	requisition := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 42585,
		IP:   net.ParseIP("0.0.0.0"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Print(err)
		return
	}
	for {
		_, remoteaddr, err := ser.ReadFromUDP(requisition)
		if err != nil {
			fmt.Print(err)
			continue
		}
		SRHUdp(ser, remoteaddr, requisition)
	}
}
