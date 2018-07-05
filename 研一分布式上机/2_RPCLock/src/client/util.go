package main

import (
	"../thrift/example"
	"log"
)

func MultiDoLock(num int64){
	//rspList := make([]*example.Rsp,num)
	for i:=0;i<int(num);i++{				//FIXME  仍然不是并行请求  在server端进行random sleep
											//solution: client的连接池  而不是routine
		req := example.Req{CliID:int64(i),Operator:"acquire"}
		rsp, err := RPCClient.DoLock(&req); if err != nil {
			log.Print(err, rsp)
			return
		}
		log.Print( "ClockID:",rsp.CliID," State:",rsp.Operator)
	}
	////由于rsp返回都是相同的
	//time.Sleep(1 * time.Second)			//wait server rsp
	//log.Print("State:",rspList[0].Operator," clientID:",rspList[0].CliID)
}

//只需要解锁原来有锁的client
func MultiUnLock(num int64){
	rspList := make([]*example.Rsp,num)
	for i:=0;i<int(num);i++{
		req := example.Req{CliID:num,Operator:"release"}
		rsp, err := RPCClient.UnLock(&req); if err != nil{
			rspList[i] = rsp
		}
	}
	//由于rsp返回都是相同的
	log.Print("State:",rspList[0].Operator," clientID:",rspList[0].CliID)
}


func UnLockOne(ID int64){
	req := example.Req{CliID:ID,Operator:"release"}
	rsp, err := RPCClient.UnLock(&req);if err != nil{
		log.Fatal(err)
		return
	}
	//由于rsp返回都是相同的
	log.Print("State:",rsp.Operator," clientID:",rsp.CliID)
}