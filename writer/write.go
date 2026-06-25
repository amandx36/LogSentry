package writer

import (
	"LogSentry/models"
	"bufio"
	"fmt"
	"os"
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

	// Write buffered data to disk
	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

func OutputWriting(report models.LogReport) error {

	if err := writeLogs("logs/output/Error.log", report.Errors); err != nil {
		return err
	}

	if err := writeLogs("logs/output/Warn.log", report.Warns); err != nil {
		return err
	}

	if err := writeLogs("logs/output/Info.log", report.Infos); err != nil {
		return err
	}

	if err := writeLogs("logs/output/Unknown.log", report.Unknown); err != nil {
		return err
	}

	fmt.Println("FILE LOGS GenEratED")

	return nil
}