package parser

import (
	"LogSentry/internal/metrics"
	"LogSentry/internal/models"
	"bufio"
	"os"
	"regexp"
	"time"
)

// best for single file parsing and loading into the buffer

// func ParseSingleFile(file *os.File, myLogs *models.LogReport) error {

// 	scanner := bufio.NewScanner(file)

// 	reg := regexp.MustCompile(`^(\S+\s+\S+)\s+(\w+)\s+\[(.*?)\]\s+(.*)$`)

// 	for scanner.Scan() {

// 		line := scanner.Text()

// 		match := reg.FindStringSubmatch(line)

// 		if match == nil {
// 			continue
// 		}

// 		entry := models.LogEntry{
// 			TimeStamp: match[1],
// 			Category:  match[2],
// 			Source:    match[3],
// 			Details:   match[4],
// 		}

// 		switch entry.Category {

// 		case "ERROR":
// 			myLogs.Errors = append(myLogs.Errors, entry)
// 			myLogs.Counts["ERROR"]++

// 		case "WARN":
// 			myLogs.Warns = append(myLogs.Warns, entry)
// 			myLogs.Counts["WARN"]++

// 		case "INFO":
// 			myLogs.Infos = append(myLogs.Infos, entry)
// 			myLogs.Counts["INFO"]++

// 		default:
// 			myLogs.Unknown = append(myLogs.Unknown, entry)
// 			myLogs.Counts["DEFAULT"]++
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return err
// 	}

// 	return nil
// }

func ParseSingleFile(file *os.File) (models.LogReport, error) {

	start := time.Now()
	defer metrics.ParseDuration.Observe(time.Since(start).Seconds())

	// Create LogReport ->Parse -> Fill Report -> Return Report
	myLogs := models.LogReport{
		Counts: models.Counts{
			"ERROR":   0,
			"WARN":    0,
			"INFO":    0,
			"DEFAULT": 0,
		},
	}

	scanner := bufio.NewScanner(file)

	reg := regexp.MustCompile(`^(\S+\s+\S+)\s+(\w+)\s+\[(.*?)\]\s+(.*)$`)

	for scanner.Scan() {

		line := scanner.Text()

		match := reg.FindStringSubmatch(line)

		if match == nil {
			metrics.ParserFailures.Inc()
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
			recordParsedLogMetrics(entry.Category)

		case "WARN":
			myLogs.Warns = append(myLogs.Warns, entry)
			myLogs.Counts["WARN"]++
			recordParsedLogMetrics(entry.Category)

		case "INFO":
			myLogs.Infos = append(myLogs.Infos, entry)
			myLogs.Counts["INFO"]++
			recordParsedLogMetrics(entry.Category)

		default:
			myLogs.Unknown = append(myLogs.Unknown, entry)
			myLogs.Counts["DEFAULT"]++
			recordParsedLogMetrics(entry.Category)
		}
	}

	if err := scanner.Err(); err != nil {
		metrics.ParserFailures.Inc()
		return models.LogReport{}, err
	}

	return myLogs, nil

}

func recordParsedLogMetrics(category string) {
	metrics.LogsProcessed.Inc()

	switch category {
	case "ERROR":
		metrics.ErrorLogs.Inc()
	case "WARN":
		metrics.WarnLogs.Inc()
	case "INFO":
		metrics.InfoLogs.Inc()
	default:
		metrics.UnknownLogs.Inc()
	}
}
