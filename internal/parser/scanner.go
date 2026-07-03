package parser

import (
	"LogSentry/internal/config"
	"os"
	"strings"
)
func ScanAndSendJobs(cfg config.Config, jobs chan<- Job) error {

// 	Directory->ReadDir() -> Send Jobs

	files, err := os.ReadDir(cfg.InputDir)
	if err != nil {
		return err
	}

	for _, value := range files {

		if value.IsDir() {
			continue
		}

		if !strings.HasSuffix(value.Name(), ".log") {
			continue
		}

		jobs <- Job{FilePath: cfg.InputDir + "/" + value.Name()}
	}
	close(jobs)
	return nil
}