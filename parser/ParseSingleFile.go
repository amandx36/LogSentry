package parser
import (
	"os"
	"bufio"
	"LogSentry/models"
	"strings"
)

func ParseSingleFile(file *os.File, myLogs *models.LogReport) error {

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		line := scanner.Text()

		parts := strings.Fields(line)

		if len(parts) < 5 {
			continue
		}

		entry := models.LogEntry{
			TimeStamp: parts[0] + " " + parts[1],
			Category:  parts[2],
			Source:    strings.Trim(parts[3], "[]"),
			Details:   strings.Join(parts[4:], " "),
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