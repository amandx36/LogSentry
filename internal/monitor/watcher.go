package monitor

import (
	"database/sql"
	"log"

	"LogSentry/internal/config"
	"LogSentry/internal/metrics"

	"github.com/fsnotify/fsnotify"
)

func DirWatching(cfg config.Config, db *sql.DB) {
	offsetManager := NewOffsetManager()

	// has 2 channel
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	// 	                Watcher
	//                    │
	//       ┌────────────┴────────────┐
	//       ▼                         ▼
	//  Events Channel           Errors Channel

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("Event:", event)
				metrics.WatcherEvents.Inc()
				if event.Op&fsnotify.Write == fsnotify.Write {
					file := event.Name

					metrics.LiveEvents.Inc()

					lastOffset := offsetManager.GetOffset(file)

					data, newOffset, err := ReadNewLogs(file, lastOffset)
					if err != nil {
						metrics.ReadFailures.Inc()
						log.Println("Error reading new logs:", err)
						continue
					}
					offsetManager.UpdateOffset(file, newOffset)
					// now send the data for processing
					log.Println("Sending data to ProcessLogs")
					ProcessLogs(data, db)
					log.Println("Sending data to ProcessLogs")
					log.Println("Modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add(cfg.InputDir)
	if err != nil {
		log.Fatal(err)
	}

	<-done
	log.Println("Sending data to ProcessLogs")
}
