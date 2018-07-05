package main

import (
	"log"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bufio"
	"os"
	"time"
)

type Config struct {
	ID           string
	LocalAddress string    `json:"localAddress"`
	Channels     []Channel `json:"channels"`
}

type Channel struct {
	RemoteAddress string `json:"remoteAddress"`
	ChannelID     string `json:"channelID"`
}

func getConnFromFile(filename string) (*MyConnection,error) {
	outputC := make(chan string)
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return nil,err
	}
	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err)
		return nil,err
	}
	conn := NewConnection(config.ID, config.LocalAddress, config.Channels, &outputC)
	fmt.Print(conn.id," ",conn.amount," ",conn.localAddr," ",conn.record," ",conn.remoteAddr[0].RemoteAddress,"\n")
	return conn,nil
}

func main() {
	process1,err := getConnFromFile("test1.json");if err != nil{
		log.Fatal(err)
		return
	}
	process2,err2 := getConnFromFile("test2.json");if err2 != nil{
		log.Fatal(err2)
		return
	}

	go func() {
		process1.msgLoop()
	}()

	go func() {
		process2.msgLoop()
	}()

	time.Sleep(2 * time.Second)
	process1.SendMsg()					//wait listen && 启动msg的轮转
	fmt.Print("输入字符m==start snapshot 输入字符c==clear record\n")

	for{
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if text[:len(text)-1] == "m"{
			process1.SendMarker()
		}
		if text[:len(text)-1] == "c"{
			process1.record = false
			process2.record = false
		}
	}
	return
}
