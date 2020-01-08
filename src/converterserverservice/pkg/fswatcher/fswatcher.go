package fswatcher

import (
	"log"
	"github.com/fsnotify/fsnotify"

)
// follow that example :
// https://gist.github.com/hiroakis/8968d8caae55d460a80206ea4930bab7

func Setup() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
	    log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
	    for {
	        select {
	        case event, ok := <-watcher.Events:
	            if !ok {
	                return
	            }
	            log.Println("event:", event)
	            if event.Op&fsnotify.Write == fsnotify.Write {
	                log.Println("modified file:", event.Name)
	            }
	        case err, ok := <-watcher.Errors:
	            if !ok {
	                return
	            }
	            log.Println("error:", err)
	        }
	    }
	}()

	err = watcher.Add(config.getTmpFolder())
	if err != nil {
	    log.Fatal(err)
	}
	<-done
}