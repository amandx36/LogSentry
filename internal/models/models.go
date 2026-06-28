package models
import (
	"time"
)

type DashBoardDetails struct{
	TotalLogs int
    Errors    int
    Warns     int
    Infos     int
    Unknown   int
}
type LogEntry struct{
	TimeStamp string 
	Category string 
	Source string 
	Details string   
}

type Counts = map[string]int
// {
// 		"ERROR":0,
// 		"WARN":0,
// 		"INFO":0,
// 		"DEFAULT":0,
// 	}

	type LogReport struct {
    Errors  []LogEntry
    Warns   []LogEntry
    Infos   []LogEntry
    Unknown []LogEntry

    Counts Counts

}

type MonitringConfig struct {

	InputFile string 
	OuptFile string 
	LiveMode bool 
	UserDetails string // for future is this authorized or not 

}

// for analytics 
type Analytics struct {
    TotalLogs      int
    TotalErrors    int
    TotalWarns     int
    TotalInfos     int
    TotalUnknown   int

    ErrorRate      float64
    WarningRate    float64

    TopSource      string
    TopSourceCount int

    GeneratedAt    time.Time
}
