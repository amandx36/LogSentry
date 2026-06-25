package parser

import (
	"LogSentry/models"
	"bufio"
	"fmt"
	"os"
	"strings"
	
)

func LoadingBuffer() (models.LogReport,models.DashBoardDetails,error){
	myLogs := models.LogReport{
		Counts: models.Counts{
			"ERROR":   0,
			"WARN":    0,
			"INFO":    0,
			"DEFAULT": 0,
		},
	}
		myDash := models.DashBoardDetails{}

	file, err := os.Open("logs/inputs/main.log")

	if err != nil {
		fmt.Printf("Error in opening the File: %v\n", err)
		 
		return models.LogReport{},myDash,err
	}

	defer file.Close()

	// buffer input/ output
	scanner := bufio.NewScanner(file)
	// intilize the struct
	
	

	for scanner.Scan() {

	line := scanner.Text()

	parts := strings.Fields(line)

	if len(parts) < 5 {
		continue
	}

	// Create LogEntry
	entry := models.LogEntry{
		TimeStamp: parts[0] + " " + parts[1],
		Category:  parts[2],
		Source:    strings.Trim(parts[3], "[]"),
		Details:   strings.Join(parts[4:], " "),
	}

	switch entry.Category {

	case "ERROR":
		myLogs.Errors = append(myLogs.Errors, entry)
		myLogs.Counts["ERROR"]++
		

	case "WARN":
		myLogs.Warns = append(myLogs.Warns, entry)
		myLogs.Counts["WARN"]++

	case "INFO":
		myLogs.Infos = append(myLogs.Infos, entry)
		myLogs.Counts["INFO"]++

	default:
		myLogs.Unknown = append(myLogs.Unknown, entry)
		myLogs.Counts["DEFAULT"]++
		}
	}
	myDash.Errors = myLogs.Counts["ERROR"]
	myDash.Infos = myLogs.Counts["INFO"]
	myDash.Warns = myLogs.Counts["WARN"]
	myDash.Unknown = myLogs.Counts["DEFAULT"]
	myDash.TotalLogs = myDash.Errors + myDash.Warns + myDash.Infos + myDash.Unknown
	
	fmt.Println(myDash)
	
	// fmt.Println("Errors")
	// for _, log := range myLogs.Errors {
	// fmt.Printf("%+v\n", log)
	// }
	// fmt.Println("INFO")
	// for _, log := range myLogs.Infos{
	// fmt.Printf("%+v\n", log)
	// }
	// fmt.Println("WARN")
	// for _, log := range myLogs.Warns {
	// fmt.Printf("%+v\n", log)
	// }
	// fmt.Println("UNKNOWN")
	// for _, log := range myLogs.Unknown {
	// fmt.Printf("%+v\n", log)
	
	// }

	if wrong := scanner.Err(); wrong != nil {
		fmt.Print("Error while loading the Scanner\n", wrong)
		return models.LogReport{} , myDash,wrong

	}
	
	fmt.Println("The Object i got ")
	fmt.Println(myLogs)
	 return myLogs, myDash,nil

	}

	

