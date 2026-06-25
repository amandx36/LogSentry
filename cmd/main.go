package main

import (
    "fmt"
    "LogSentry/parser"
)

func main() {
    report, err := parser.MainParsing()
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(report)
}