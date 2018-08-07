package main

import (
	"fmt"
	"time"
)

func main() {
	nTime := time.Now()
	//yesTime := nTime.AddDate(0, 0, 0)
	//logDay := yesTime.Format("20060102 00:00:00")
	logDay := nTime.Format("20060102 00:00:00")
	fmt.Println(logDay)
}
