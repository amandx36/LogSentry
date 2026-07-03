package parser

import (
	"LogSentry/internal/config"
	"LogSentry/internal/models"
	"runtime"
	
)

// Best for single file parsing and loading into the buffer

// func LoadingBuffer(cfg config.Config) (models.LogReport, models.DashBoardDetails, error) {

// 	myLogs := models.LogReport{
// 		Counts: models.Counts{
// 			"ERROR":   0,
// 			"WARN":    0,
// 			"INFO":    0,
// 			"DEFAULT": 0,
// 		},
// 	}

// 	myDash := models.DashBoardDetails{}

// 	files, err := os.ReadDir(cfg.InputDir)
// 	if err != nil {
// 		return models.LogReport{}, myDash, err
// 	}

// 	for _, value := range files {

// 		if value.IsDir() {
// 			continue
// 		}

// 		if !strings.HasSuffix(value.Name(), ".log") {
// 			continue
// 		}

// 		file, err := os.Open(cfg.InputDir + "/" + value.Name())
// 		if err != nil {
// 			continue
// 		}

// 		err = ParseSingleFile(file, &myLogs)
// 		file.Close()

// 		if err != nil {
// 			return models.LogReport{}, myDash, err
// 		}
// 	}

// 	myDash.Errors = myLogs.Counts["ERROR"]
// 	myDash.Warns = myLogs.Counts["WARN"]
// 	myDash.Infos = myLogs.Counts["INFO"]
// 	myDash.Unknown = myLogs.Counts["DEFAULT"]
// 	myDash.TotalLogs = myDash.Errors +
// 		myDash.Warns +
// 		myDash.Infos +
// 		myDash.Unknown

// 	return myLogs, myDash, nil
// }


// for concurrent parsing and loading into the buffer
// WorkFlow 

// Create Jobs Channel -> create Results Channel -> Start Workers -> Scanner sends Jobs -> Workers Parse Files -> Workers send Results -> MergeReports() -> Final Report

func LoadingBuffer(cfg config.Config) (models.LogReport, models.DashBoardDetails, error){
	
	

numWorkers := runtime.NumCPU()  // return number of cpu core dude 
	// buffer channel  for job holding 
	jobs := make(chan Job, numWorkers)
	// buffer channel for results holding 
	results := make(chan Result, numWorkers)

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go Worker(w, jobs, results)
	}

	// Scan and send jobs
	jobCount, err := ScanAndSendJobs(cfg, jobs)
	if err != nil {
		return models.LogReport{}, models.DashBoardDetails{}, err
	}


	finalReport := models.LogReport{
		Counts: models.Counts{
			"ERROR":   0,
			"WARN":    0,
			"INFO":    0,
			"DEFAULT": 0,
		},
	}

	for i := 0; i < jobCount; i++ {
		result := <-results
		if result.Err != nil {
			continue 
		}
		MergeReports(&finalReport, result.Report)
	}

	myDash := models.DashBoardDetails{
		Errors:  finalReport.Counts["ERROR"],
		Warns:   finalReport.Counts["WARN"],
		Infos:   finalReport.Counts["INFO"],
		Unknown: finalReport.Counts["DEFAULT"],
		TotalLogs: finalReport.Counts["ERROR"] +
			finalReport.Counts["WARN"] +
			finalReport.Counts["INFO"] +
			finalReport.Counts["DEFAULT"],
	}

	return finalReport, myDash, nil

}