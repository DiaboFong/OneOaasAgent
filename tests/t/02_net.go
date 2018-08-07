package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.DialTimeout("tcp", "centos7:8180", 1*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("端口已开启")
	defer conn.Close()
}
