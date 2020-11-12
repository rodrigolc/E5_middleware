package main

import (
	"fmt"
	"net/rpc"
	"encoding/json"
	"net"
	"bufio"
	"strconv"
	"strings"
)


func sendRecieveMessage(connection int,operation int, candidate int) {
	var replyTCP string
	request := struct { 
		Op int
		Ca int 
	}{ 
		Op: operation, 
		Ca: candidate, 
	}
	args, _ := json.Marshal(request)
	if connection == 1 {
		tcpConnection, err := rpc.Dial("tcp", "localhost:42586")
		if err != nil {
			fmt.Print(err)
			return 
		}
		
		err = tcpConnection.Call("Listener.SRHTcp", []byte(args), &replyTCP)
		if err != nil {
			fmt.Print(err)
			return
		}
		fmt.Println(replyTCP)
	} else {
		buffer :=  make([]byte, 2048)
		udpConnection, err := net.Dial("udp", "localhost:42585")
		if err != nil {
			fmt.Print(err)
			return
		}
		msg :=  strconv.Itoa(operation) +":"+ strconv.Itoa(candidate)+";" 
		fmt.Fprintf(udpConnection, msg)
		_, err = bufio.NewReader(udpConnection).Read(buffer)
		fmt.Print(strings.Split(string(buffer), ";")[0])		
    udpConnection.Close()
	}
}

func Client() {
	var operation int
	var connection int 
	var candidate int
	for true {
		fmt.Println("Qual tipo de conexão você deseja utilizar?\n\n 1 - TCP\n 2 - UDP\n")
		fmt.Scanln(&connection)
	
		fmt.Println("Você deseja registrar um voto ou ver a parcial?\n\n 1 - Registrar Voto \n 2 - Ver Parcial\n")
		fmt.Scanln(&operation)

		if operation == 1 {
			fmt.Println("Escolha um candidato?\n\n 1 - A \n 2 - B\n 3 - C\n")
			fmt.Scanln(&candidate)
			sendRecieveMessage(connection, operation, candidate)
		} else if operation == 2 {
			sendRecieveMessage(connection, operation, 0)
		} else {
			fmt.Printf("Operação inválida")
		}
	}
}
