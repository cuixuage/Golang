package main

import (
	"../thrift/example"
	"git.apache.org/thrift.git/lib/go/thrift"
	"fmt"
	"log"
	"sync"
	"encoding/json"
)

const (
	HOST = "localhost"
	PORT = "6790"
)

var MapRWMutex *sync.RWMutex
//维护map保存lock acquire,操作时使用读写锁保护
var ClientState map[int64]bool = make(map[int64]bool,0)

//实现handler *GetLock的接口
type FormatDataImpl struct {}

func (fdi *FormatDataImpl) DoLock(req *example.Req) (*example.Rsp,error){
	MapAddNewClient(req.CliID)
	RandomLock()
	lockID := GetCurrentLock()
	var rsp example.Rsp
	rsp.CliID = lockID
	rsp.Operator = "acquire"
	//log.Print("clockID=",rsp.CliID," operator=",rsp.Operator)
	return &rsp, nil
}

func (fdi *FormatDataImpl) UnLock(req *example.Req) (*example.Rsp,error){
	UnlockClient(req.CliID)
	var rsp example.Rsp
	rsp.CliID = req.CliID
	rsp.Operator = "release"
	log.Print("clockID=",rsp.CliID," operator=",rsp.Operator)
	return &rsp, nil
}

func (fdi *FormatDataImpl) ClientStates() (*example.Rsp,error){
	var rsp example.Rsp
	rsp.CliID = -1
	rsp.Operator = "-1"
	jsonString, err := json.Marshal(ClientState); if err==nil{
		rsp.Buffer = string(jsonString)
	}
	return &rsp, nil
}

func StartSever(){
	handler := &FormatDataImpl{}
	processor := example.NewGetLockProcessor(handler)

	serverTransport, err := thrift.NewTServerSocket(HOST + ":" + PORT)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("Running at:", HOST + ":" + PORT)
	server.Serve()
}

func main() {
	MapRWMutex = new(sync.RWMutex)
	StartSever()
}