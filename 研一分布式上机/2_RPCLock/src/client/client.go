package main

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"net"
	"../thrift/example"
	"log"
	"bufio"
	"os"
	"strconv"
	"strings"
	"fmt"
)

const (
	HOST = "localhost"
	PORT = "6790"
)

var RPCClient *example.GetLockClient

func main() {
	tSocket, err := thrift.NewTSocket(net.JoinHostPort(HOST, PORT))
	if err != nil {
		log.Fatalln("tSocket error:", err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	transport := transportFactory.GetTransport(tSocket)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	RPCClient = example.NewGetLockClientFactory(transport, protocolFactory)
	if err := transport.Open(); err != nil {
		log.Fatalln("Error opening:", HOST+":"+PORT)
	}
	defer transport.Close()

	fmt.Print("e.g. acquire,5(不同client请求次数)\n")
	fmt.Print("e.g. release,1(relase clientID)\n")
	fmt.Print("e.g. status\n")
	for{
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]
		textList := strings.Split(text,",")
		if len(textList) == 1{
			operator := strings.Split(text,",")[0]
			operator = strings.Replace(operator, " ", "", -1)
			if operator == "status"{
				rsp,err := RPCClient.ClientStates(); if err == nil{
					log.Print(rsp.Buffer)
				}
			}
			continue
		}else if len(textList) == 2{
			operator := strings.Split(text,",")[0]
			operator = strings.Replace(operator, " ", "", -1)
			ID := strings.Split(text,",")[1]
			ID = strings.Replace(ID, " ", "", -1)
			num, err := strconv.ParseInt(ID, 10, 64);if err != nil {
				log.Fatal(err)
			}
			if operator == "acquire"{
				MultiDoLock(num)
				continue
			}
			if operator == "release" {
				UnLockOne(num)
				continue
			}
		}

		fmt.Print("input correct cmd\n")

	}
	return
}