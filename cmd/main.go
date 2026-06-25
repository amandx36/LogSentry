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

	report, _, err := parser.LoadingBuffer(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = writer.OutputWriting(cfg.OutputDir, report)
	if err != nil {
		fmt.Println(err)
	}

	
}