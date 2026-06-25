package main

import (
	"bufio"
	"fmt"
	
	"os"
	"strings"
)

func LoadingBuffer(){
		var (
    		Errors  []string
    		Warn    []string
    		Info    []string
    		Unknown []string
	)
	counts := map[string]int{
		"ERROR":0,
		"WARN":0,
		"INFO":0,
		"DEFAULT":0,
	}

	// open the file 
			file, err := os.Open("Logs/main.log")
	 		
			if err != nil {
				fmt.Printf("Error in opening the File: %v\n", err)	
				return 
			}
			defer file.Close()
			
			// buffer input/ output 
			scanner := bufio.NewScanner(file)
		
			for scanner.Scan() {
				line := scanner.Text()
    			fmt.Println(line) 
		
				switch{
					case strings.Contains(line,"ERROR"):
						Errors = append(Errors, line)
						counts["ERROR"]++
			
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
				if wrong := scanner.Err()
			wrong != nil{
				fmt.Print("Error while loading the Scanner\n",wrong)
				return 

			}
			fmt.Println(counts["ERROR"])
			fmt.Println(counts["WARN"])
			fmt.Println(counts["INFO"])
			fmt.Println(counts["DEFAULT"])
			fmt.Println("Total Errors ")
			for i := 0 ; i < len(Errors) ; i++{
				fmt.Println(Errors[i]);
			}
			fmt.Println("Total number of Warns ")
			for i := 0 ; i < len(Warn); i++{
				fmt.Println(Warn[i])
			}
			fmt.Println("Total number of Info ")
			for i := 0 ; i < len(Warn); i++{
				fmt.Println(Info[i])
			}
			
			
		}


func main (){
	LoadingBuffer();
}
