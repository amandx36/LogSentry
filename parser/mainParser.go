package parser

import (
	"LogSentry/models"
	"bufio"
	"fmt"
	"os"
	"strings"
	"regexp"
)

func LoadingBuffer() {

	file, err := os.Open("Logs/main.log")

	if err != nil {
		fmt.Printf("Error in opening the File: %v\n", err)
		return
	}
	defer file.Close()

	// buffer input/ output
	scanner := bufio.NewScanner(file)
	// intilize the struct
	myDash := models.DashBoardDetails{}
	
// 	type LogReport struct {
//     Errors  []LogEntry
//     Warns   []LogEntry
//     Infos   []LogEntry
//     Unknown []LogEntry

//     Counts Counts

// }

// type LogEntry struct{
// 	TimeStamp string 
// 	Category string 
// 	Details string   
// }





	logBio := models.LogEntry{}

	myLogs := models.LogReport{
		Occurs: models.Counts{
			"ERROR":   0,
			"WARN":    0,
			"INFO":    0,
			"DEFAULT": 0,
		},
	}
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		switch {
		case strings.Contains(line, "ERROR"):
			`2026-06-24 13:21:57.767 WARN  [devmind-api] o.s.b.a.e.web.EndpointLinksResolver : Exposure of /actuator/env not recommended in production`
			parts : strings.Split()
			
			myLogs.Errors = append(myLogs.Errors, )

		case strings.Contains(line, "WARN"):
			Warn = append(Warn, line)
			counts["WARN"]++

		case strings.Contains(line, "INFO"):
			Info = append(Info, line)
			counts["INFO"]++

		default:
			Unknown = append(Unknown, line)
			counts["DEFAULT"]++

		}

	}
	if wrong := scanner.Err(); wrong != nil {
		fmt.Print("Error while loading the Scanner\n", wrong)
		return

	}
	fmt.Println(counts["ERROR"])
	fmt.Println(counts["WARN"])
	fmt.Println(counts["INFO"])
	fmt.Println(counts["DEFAULT"])
	fmt.Println("Total Errors ")
	for i := 0; i < len(Errors); i++ {
		fmt.Println(Errors[i])
	}
	fmt.Println("Total number of Warns ")
	for i := 0; i < len(Warn); i++ {
		fmt.Println(Warn[i])
	}
	fmt.Println("Total number of Info ")
	for i := 0; i < len(Warn); i++ {
		fmt.Println(Info[i])
	}

}

func main() {
	LoadingBuffer()
}
