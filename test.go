package main

import(
	"fmt"
	"time"
)

func main(){
	fmt.Println(time.Now().String())
	fmt.Println(time.Now().String()[:4])
	fmt.Println(time.Now().String()[5:7])
	return
}