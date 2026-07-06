package parser

import (
	"LogSentry/internal/config"
	"LogSentry/internal/metrics"
	"os"
	"strings"
)

func ScanAndSendJobs(cfg config.Config, jobs chan<- Job) (int, error) {

	files, err := os.ReadDir(cfg.InputDir)
	if err != nil {
		return 0, err
	}

	count := 0

	for _, value := range files {

		if value.IsDir() {
			continue
		}

		if !strings.HasSuffix(value.Name(), ".log") {
			continue
		}

		jobs <- Job{
			FilePath: cfg.InputDir + "/" + value.Name(),
		}
		metrics.FilesProcessed.Inc()
		count++
	}

	close(jobs)

	return count, nil
}
