package fswatcher

import (
	"log"
	"errors"
	"os"
	
	fslib "rapGO.io/src/converterserverservice/pkg/fswatcher/lib"

)


func Setup() {
	watcher, err := fslib.NewWatcher()
	if err != nil {
	    log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
	    for {
	        select {
	        case event := <-watcher.Events:
	            log.Println("event:", event)
	            if event.Op&fslib.Create == fslib.Create {
	                log.Println("modified file:", event.Name)
				}
				
				if event.Op&fslib.Remove == fslib.Remove {

				}
	        case err, ok := <-watcher.Errors:
	            if !ok {
	                return
	            }
	            log.Println("error:", err)
	        }
	    }
	}()
	tmpFolder, ok := os.LookupEnv("TMP_FOLDER")
	if !ok {
		panic(errors.New("The environment variable TMP_FOLDER is not defined."))
	}

	err = watcher.Add(tmpFolder)
	if err != nil {
	    log.Fatal(err)
	}
	<-done
}