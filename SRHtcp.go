package main

import (
	"log"
	"net"
	"net/rpc"
	"encoding/json"
	"strconv"
)

type Listener int

type Args struct {
	operation int
	candidate int
}

var candidates [3]int 


func registerVote(candidate int) {
	candidates[candidate-1] += 1
}

func getVotes() string {
	msg := "\n\n A: " + strconv.Itoa(candidates[0]) + " B: " + strconv.Itoa(candidates[1]) + " C: " + strconv.Itoa(candidates[2]) + "- TCP\n\n"
	return msg
}

func TranformData(requisition []byte) Args {
	var dat map[string]interface{}
	json.Unmarshal(requisition, &dat)
	args := Args { int(dat["Op"].(float64)) , int(dat["Ca"].(float64)) }
	return args
}

func (l *Listener) SRHTcp(requisition []byte, reply *string) error {
	var args = TranformData(requisition)
	if args.operation == 1 {
		registerVote(args.candidate)
		*reply = "\n\n Seu Voto foi computado com sucesso - TCP\n\n"
	} else {
		*reply = getVotes()
	}
	return nil
}

func main() {

	addy, err := net.ResolveTCPAddr("tcp", "0.0.0.0:42586")
	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		log.Fatal(err)
	}	

	listener := new(Listener)
	rpc.Register(listener)
	rpc.Accept(inbound)
}
