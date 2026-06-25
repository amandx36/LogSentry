package main

import (
	"LogSentry/config"
	"LogSentry/parser"
	"LogSentry/writer"
	"fmt"
)

func main() {

	cfg, err := config.Loadconfig("config/config.json")
	if err != nil{
		fmt.Println("Got Error ",err)
		return 
	}
	

	
	// report, _,err := parser.LoadingBuffer()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// err = writer.OutputWriting(report)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}