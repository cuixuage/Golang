package main

import (
	"fmt"
	"net"
	"os"
)

func testIP() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0]) //args[0]是terminal空格切割第一个字符串
		os.Exit(1)
	}
	name := os.Args[1]

	//pareseIP 也能够判断IP地址有效性   net.pareseIP(IP String)
	addr := net.ParseIP(name)
	if addr == nil {
		fmt.Println("Invalid address")
	} else {
		fmt.Println("The address is ", addr.String())
	}
	os.Exit(0)
}

func main() {
	testIP()
}
