package parser

import (
	"LogSentry/internal/config"
	"testing"
)

func TestLoadingBuffer(t *testing.T) {

	cfg, err := config.Loadconfig("../config/config.json")
	if err != nil {
		t.Fatal(err)
	}

	report, _, err := LoadingBuffer(cfg)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(report)
}