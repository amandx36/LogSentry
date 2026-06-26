package main

import (
	"LogSentry/internal/config"
	"LogSentry/internal/dashboard"
	"LogSentry/internal/parser"
	"LogSentry/internal/storage/postgres"
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

	db, err := postgres.Connect(cfg)
	 if err != nil {
		fmt.Println(err)
		return
		}
	

	err = postgres.InsertLogs(db, report)
		if err != nil {
			fmt.Println(err)
		return
	}
	fmt.Println("Getting All statics from the database :) ")
	
	StatsInfo , err  := dashboard.GetOverallStats(db)
	defer db.Close()
	if err != nil {
		fmt.Println("Error while getting the info from the database :)")
		return 
	}
	fmt.Println("Total Errors " ,StatsInfo.Errors )
	fmt.Println("Total Info " ,StatsInfo.Infos )
	fmt.Println("Total Total Logs  " ,StatsInfo.TotalLogs )
	fmt.Println("Total Unknown  " ,StatsInfo.Unknown )
	fmt.Println("Total Warns  " ,StatsInfo.Warns )

}