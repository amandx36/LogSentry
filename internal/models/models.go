package models


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

