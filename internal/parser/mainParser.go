package parser

import (
	"LogSentry/internal/config"
	"LogSentry/internal/models"
	"runtime"
	"sync"
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

func LoadingBuffer(cfg config.Config) (models.LogReport, models.DashBoardDetails, error) {

	// Number of workers = Number of logical CPU cores
	numWorkers := runtime.NumCPU()

	// Channels
	jobs := make(chan Job, numWorkers)
	results := make(chan Result, numWorkers)

	// WaitGroup
	var wg sync.WaitGroup

	// Start workers
	wg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go Worker(i, jobs, results, &wg)
	}

	// Scanner sends jobs
	jobCount, err := ScanAndSendJobs(cfg, jobs)
	if err != nil {
		return models.LogReport{}, models.DashBoardDetails{}, err
	}

	// Close results after every worker exits
	go func() {
		wg.Wait()
		close(results)
	}()

	// Final Report
	finalReport := models.LogReport{
		Counts: models.Counts{
			"ERROR":   0,
			"WARN":    0,
			"INFO":    0,
			"DEFAULT": 0,
		},
	}

	// Receive all parsed reports
	for i := 0; i < jobCount; i++ {

		result := <-results

		if result.Err != nil {
			continue
		}

		MergeReports(&finalReport, result.Report)
	}

	dashboard := models.DashBoardDetails{
		Errors:  finalReport.Counts["ERROR"],
		Warns:   finalReport.Counts["WARN"],
		Infos:   finalReport.Counts["INFO"],
		Unknown: finalReport.Counts["DEFAULT"],
	}

	dashboard.TotalLogs =
		dashboard.Errors +
			dashboard.Warns +
			dashboard.Infos +
			dashboard.Unknown

	return finalReport, dashboard, nil
}
