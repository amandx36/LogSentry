package main

import (
	"LogSentry/internal/config"
	"LogSentry/internal/dashboard"
	"LogSentry/internal/parser"
	"LogSentry/internal/search"
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

	searchCat, err := search.SearchByCategory(db, "ERROR")
	if err != nil{
		return 
	}else{
			fmt.Println("The Error i got dude :) ",searchCat)

	}
	

	Search , err := search.SearchBySource(db,"redish")
	if err != nil{
		fmt.Println("Error in Searching  by Source") 
	}else{
		fmt.Println("The Search u got ",Search)

	}
	
	keywordLogs, err := search.SearchByKeywords(db, "timeout")
	if err != nil {
		fmt.Println("Error searching by keyword:", err)
	} else {
		fmt.Println("KEYWORD SEARCH: 'timeout'", keywordLogs)
	}

	recentLogs, err := search.GetRecentLogs(db, 5)
	if err != nil {
		fmt.Println("Error fetching recent logs:", err)
	} else {
		fmt.Println("TOP 5 RECENT LOGS", recentLogs)
	}
	
	


}