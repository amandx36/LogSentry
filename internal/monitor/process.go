package monitor

import (
	"LogSentry/internal/repository/postgres"
	"database/sql"
	"log"
	"LogSentry/internal/parser"
)

// data -> parser -> logReport -> batch insert


func ProcessLogs(data []byte, db *sql.DB)(error){
	 report ,  _,err :=parser.LoadingBuffer()
	 if err !=nil{
		log.Print("Error  while parsing the file")
		 
	 }
	 // insert the data 
	 if err := postgres.InsertLogs(db, report); err != nil {
    return err
}
}