package parser

import (
	"LogSentry/internal/config"
	"LogSentry/internal/models"
	"os"
	"strings"
)

func LoadingBuffer(cfg config.Config) (models.LogReport, models.DashBoardDetails, error) {

	myLogs := models.LogReport{
		Counts: models.Counts{
			"ERROR":   0,
			"WARN":    0,
			"INFO":    0,
			"DEFAULT": 0,
		},
	}

	myDash := models.DashBoardDetails{}

	files, err := os.ReadDir(cfg.InputDir)
	if err != nil {
		return models.LogReport{}, myDash, err
	}

	for _, value := range files {

		if value.IsDir() {
			continue
		}

		if !strings.HasSuffix(value.Name(), ".log") {
			continue
		}

		file, err := os.Open(cfg.InputDir + "/" + value.Name())
		if err != nil {
			continue
		}

		err = ParseSingleFile(file, &myLogs)
		file.Close()

		if err != nil {
			return models.LogReport{}, myDash, err
		}
	}

	myDash.Errors = myLogs.Counts["ERROR"]
	myDash.Warns = myLogs.Counts["WARN"]
	myDash.Infos = myLogs.Counts["INFO"]
	myDash.Unknown = myLogs.Counts["DEFAULT"]
	myDash.TotalLogs = myDash.Errors +
		myDash.Warns +
		myDash.Infos +
		myDash.Unknown

	return myLogs, myDash, nil
}