package writer

// LogReport
//       │
//       ▼
// Errors Slice
//       │
//       ▼
// Logs/Error.log

// Warn Slice
//       │
//       ▼
// Logs/Warn.log

// Infos Slice
//       │
//       ▼
// Logs/Info.log

// Unknown Slice
//       │
//       ▼
// Logs/Unknown.log
import (
	"fmt"
	"os"
	"LogSentry/parser"
)

func outputWriting (){
	path := "./logs/output"

	outPutfile , err := os.Create(path)
	if err != nil {
		fmt.Println("Error in creating the path")

	}
	defer outPutfile.Close()
	    data , err  := parser.MainParsing()
		if err != nil {
			fmt.Println("Error in Parsing #02 ")
		}
	fmt.Println("Total Number of Error ")
	
	
}
