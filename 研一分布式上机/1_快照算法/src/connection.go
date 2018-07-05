package main

import (
	"net"
	"io"
	"log"
	"time"
	"strings"
)

type MyConnection struct {
	amount     int
	id         string
	localAddr  string
	remoteAddr []Channel
	outputC    *chan string				//用来记录routine中的information
	record     bool
}

func NewConnection(id, laddr string, raddr []Channel, outputC *chan string) *MyConnection {
	return &MyConnection{
		amount:     101,
		id:         id,
		localAddr:  laddr,
		remoteAddr: raddr,
		outputC:    outputC,
		record:     false,
	}
}

func (th *MyConnection) getAmount(){
	log.Print("store snapshot in ",th.localAddr," amount=",th.amount)
}

// connects the specified :port ,establish a connection
//tcp 链接是双向的； 用record标识当作单向使用
//向所有remote_conn
func (th *MyConnection) SendMarker() {
	th.record = true
	//*th.outputC <- "START RECORD"
	//time.Sleep(2 * time.Second)
	//log.Print(th.remoteAddr," send marker")
	th.getAmount()
	for _, adress := range th.remoteAddr {
		conn, err := net.Dial("tcp", adress.RemoteAddress)
		defer conn.Close()
		if err != nil {
			log.Print(err)
			*th.outputC <- "ERROR: " + err.Error()
			return
		}
		_, err = conn.Write([]byte(th.id + "|marker"))					//发送mark标识
		if err != nil {
			log.Print(err)
			*th.outputC <- "ERROR: " + err.Error()
			return
		}
	}
}

//前提: local remote port都已经被监听listen
func (th *MyConnection) SendMsg(){
	//log.Print(th.localAddr," send Msg")
	for _, adress := range th.remoteAddr {
		conn, err := net.Dial("tcp", adress.RemoteAddress)
		if err != nil {
			log.Print("get conn fail, ",err)
			*th.outputC <- "ERROR: " + err.Error()
			return
		}
		_, err = conn.Write([]byte(th.id + "|message in transporting"))					//发送mark标识
		if err != nil {
			log.Print(err)
			*th.outputC <- "ERROR: " + err.Error()
			return
		}
	}
}

//收到轮转的消息Msg 则返回true
func (th *MyConnection) receiveMsg(listener net.Listener) (bool,error){
	_, buffer, err := th.myReceive(listener)
	if err != nil {
		*th.outputC <- "ERROR: " + err.Error()
		return false,err
	}
	//*th.outputC <- string(buffer)
	split := strings.SplitN(string(buffer), "|", -1)
	if len(split) < 2 {
		return false,err
	}
	//channelId := split[0]
	message := split[1]
	if message == "message in transporting" {
		if err != nil {
			*th.outputC <- "ERROR: " + err.Error()
			return  false,err
		}
		th.amount +=1
		//log.Print(th.localAddr," receiveMsg, count=",th.amount,"\n\n")
		return true,nil
	}
	if message == "marker" && !th.record {
		th.SendMarker()
		return false,nil
	}
	return false,nil
}

func (th *MyConnection) msgLoop(){
	listener, err := net.Listen("tcp", th.localAddr)
	if err != nil {
		*th.outputC <- "ERROR: " + err.Error()
		log.Println(err.Error())
		return
	}
	for {
		msg,err := th.receiveMsg(listener); if err != nil{
			continue
		}
		if msg == true{
			time.Sleep(2 * time.Second)
			th.SendMsg()
		}
	}
}

func (th *MyConnection) myReceive(listener net.Listener) (string, []byte, error) {
	conn, err := listener.Accept()
	if err != nil {
		*th.outputC <- "ERROR: " + err.Error()
		return "", []byte{}, err
	}
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	//log.Print("conn read ok\n")
	if err != nil && err != io.EOF {
		*th.outputC <- "ERROR: " + err.Error()
		return "", []byte{}, err
	}
	//log.Print( "msg轮转ing...")
	return conn.RemoteAddr().String(), buffer[:n], nil
}
