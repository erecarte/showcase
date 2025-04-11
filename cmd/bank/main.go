package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

const listenDirectory = "data/bank/"

// main
func main() {
	var processedFiles = map[string]bool{}
	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				if event.Op != fsnotify.Create || !strings.HasSuffix(event.Name, ".request.xml") {
					continue
				}
				b := processedFiles[event.Name]
				if b {
					continue
				}
				fmt.Printf("EVENT! %#v\n", event)
				fileID := strings.TrimPrefix(strings.TrimSuffix(event.Name, ".request.xml"), listenDirectory)
				file, err := os.Create(fmt.Sprintf("%s/%s.response.csv", listenDirectory, fileID))
				if err != nil {
					fmt.Println("ERROR", err)
				}
				defer file.Close()
				w := csv.NewWriter(file)
				defer w.Flush()
				err = w.WriteAll([][]string{
					{"ID", "STATUS"},
					{fileID, "PROCESSED"},
				})
				if err != nil {
					fmt.Println(fmt.Sprintf("ERROR processing fileID=%s:", fileID), err)
				}
				processedFiles[event.Name] = true
				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add(listenDirectory); err != nil {
		fmt.Println("ERROR", err)
	}

	<-done
}
