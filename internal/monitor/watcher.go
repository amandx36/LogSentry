package monitor

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"LogSentry/internal/config"
)

func StartMonitoring(cfg config.Config){
	

	// has 2 channel 
	watcher , err := fsnotify.NewWatcher()
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
				if event.Op&fsnotify.Write == fsnotify.Write {
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
}