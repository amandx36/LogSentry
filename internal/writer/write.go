package writer

import (
	"LogSentry/internal/models"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func writeLogs(path string, logs []models.LogEntry) error {

	outPutfile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outPutfile.Close()

	writer := bufio.NewWriter(outPutfile)

	for _, log := range logs {

		_, err := fmt.Fprintf(
			writer,
			"%s %s [%s] %s\n",
			log.TimeStamp,
			log.Category,
			log.Source,
			log.Details,
		)

		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func OutputWriting(outputDir string, report models.LogReport) error {

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	if err := writeLogs(filepath.Join(outputDir, "Error.log"), report.Errors); err != nil {
		return err
	}

	if err := writeLogs(filepath.Join(outputDir, "Warn.log"), report.Warns); err != nil {
		return err
	}

	if err := writeLogs(filepath.Join(outputDir, "Info.log"), report.Infos); err != nil {
		return err
	}

	if err := writeLogs(filepath.Join(outputDir, "Unknown.log"), report.Unknown); err != nil {
		return err
	}

	fmt.Println("FILE LOGS GENERATED")

	return nil
}
