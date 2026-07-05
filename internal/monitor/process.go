package monitor

import (
	"LogSentry/internal/repository/postgres"
	"database/sql"
	"log"
	"LogSentry/internal/parser"
)

// data -> parser -> logReport -> batch insert


func ProcessLogs(data []byte, db *sql.DB)(error){
	log.Println("ProcessLogs logs this for file check ")
	 report , err := parser.ParseByte(data)
	 if err !=nil{
		log.Print("Error  while parsing the file")
		return err 
		 
	 }
	 log.Printf("Parsed: Errors=%d Warns=%d Infos=%d Unknown=%d",
    report.Counts["ERROR"],
    report.Counts["WARN"],
    report.Counts["INFO"],
    report.Counts["DEFAULT"],
)
	 // insert the data 
	 if err := postgres.InsertLogs(db, report); err != nil {
    return err
	}
	log.Println("Inserted into PostgreSQL") 
	return nil ;
}