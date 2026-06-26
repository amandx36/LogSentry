package parser

import (
	"LogSentry/internal/models"
	"bufio"
	"os"
	"regexp"
)

func ParseSingleFile(file *os.File, myLogs *models.LogReport) error {

	scanner := bufio.NewScanner(file)

	reg := regexp.MustCompile(`^(\S+\s+\S+)\s+(\w+)\s+\[(.*?)\]\s+(.*)$`)

	for scanner.Scan() {

		line := scanner.Text()

		match := reg.FindStringSubmatch(line)

		if match == nil {
			continue
		}

		entry := models.LogEntry{
			TimeStamp: match[1],
			Category:  match[2],
			Source:    match[3],
			Details:   match[4],
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

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}