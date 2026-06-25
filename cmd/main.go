package main

import (
	"LogSentry/parser"
	"LogSentry/writer"
	"fmt"
)

func main() {
	report, err := parser.LoadingBuffer()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = writer.OutputWriting(report)
	if err != nil {
		fmt.Println(err)
	}
}