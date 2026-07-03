package parser

import (
	"LogSentry/internal/models"
	"os"
    "sync"
)

type Job struct {
	FilePath string
}

type Result struct {
	Report models.LogReport
	Err    error
}



// jobs <-chan Job      // Receive-only channel (reads from jobs)
// results chan<- Result // Send-only channel (writes to results)

func Worker(id int, jobs <-chan Job, results chan<- Result , wg *sync.WaitGroup)  {
    wg.Done() 
	for job := range jobs {

		// Open the file
		file, err := os.Open(job.FilePath)
		if err != nil {
			results <- Result{Err: err}
			continue
		}

		// Parse the file
		report, err := ParseSingleFile(file)

		// Close the file
		file.Close()

		if err != nil {
			results <- Result{Err: err}
			continue
		}

		// Send parsed report back
		results <- Result{
			Report: report,
			Err:    nil,
		}
	}
}
