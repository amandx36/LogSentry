package main

import (
	"LogSentry/internal/config"
	"LogSentry/internal/parser"
	"LogSentry/internal/writer"
	"fmt"
	"LogSentry/internal/storage/postgres"
	
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

	db, err := postgres.Connect(cfg)
    if err != nil {
        panic(err) // Crash if the DB is offline
    }
    
    // THIS is where the defer belongs! 
    // It keeps the DB alive until the entire program finishes.
    defer db.Close()



	
}