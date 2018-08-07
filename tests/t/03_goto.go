package main

import (
	"fmt"
)

func main() {
	var a1 = [3]int{1, 2, 3}
	for k, v := range a1 {
		fmt.Println(k, v)
		if k == 1 {
			fmt.Println("----------------------")
			goto EOF
		}
	EOF:
		fmt.Println("============")
	}
}
