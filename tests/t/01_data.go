package main

import (
	"encoding/json"
	"fmt"
)

type Data struct {
	Count int
	Data  []IpMap
}

type IpMap struct {
	IP   string
	Info InfoDetail
}
type InfoDetail struct {
	Port int
	User string
	Pwd  string
}

func main() {
	var data Data
	var ipmap IpMap
	ipmap.IP = "127.0.0.1"
	ipmap.Info.Port = 22
	ipmap.Info.User = "root"
	ipmap.Info.Pwd = "admin"
	data.Data = append(data.Data, ipmap)

	ipmap.IP = "127.0.0.1"
	ipmap.Info.Port = 22
	ipmap.Info.User = "root"
	ipmap.Info.Pwd = "admin"
	data.Data = append(data.Data, ipmap)

	data.Count = len(data.Data)
	fmt.Println(data)
	s, _ := json.Marshal(data)
	fmt.Println(string(s))
}
