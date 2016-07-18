package main

import (
	"log"
	"time"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


func main() {
	log.Printf("hello, world\n")

	watcher, err := fsnotify.NewWatcher()
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
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	go func() {
		time.Sleep(60 * time.Second)
		done <- true
	}()

	err = watcher.Add("c:\\foo")
	if err != nil {
		log.Fatal(err)
	}
	<-done


	log.Println("exiting.")

}
