package main

import (
	"LogSentry/internal/config"
	"LogSentry/internal/parser"
	"LogSentry/internal/writer"
	"fmt"
)

func main() {

	cfg, err := config.Loadconfig("internal/config/config.json")
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