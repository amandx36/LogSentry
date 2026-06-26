package postgres

import (
	"database/sql"
	"fmt"

	"golang.org/x/tools/go/analysis/passes/defers"
)

func  Connect() {
	ConnString := "postgres://postgres:root@localhost:5432/logsentry?sslmode=disable"


	// opening the database connection 
	db , err := sql.Open("postgres",ConnString)
	if err != nil{
		fmt.Println("Error in connection the Data base #1 ",err)
		return 

	}
	defer db.Close()
	// check is connection working or not
	err = db.Ping()
	if err != nil{
		fmt.Println("Failed to Connect to the database ",err)
	}
	fmt.Println("Connected successfully to the post-gre database ")


}