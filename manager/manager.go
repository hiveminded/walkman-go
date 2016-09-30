// manager is a service that will update elastic database with music directory file change
package main

import (
	"log"

	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

var (
	watcher *fsnotify.Watcher
)

func walk(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		return nil
	}

	fmt.Printf("Visited: %s\n", path)

	return watcher.Add(path)
}

func do(event fsnotify.Event) {
	log.Println("event:", event)
	if event.Op&fsnotify.Write == fsnotify.Write {
		log.Println("modified file:", event.Name)
	}
}

func main() {
	done := make(chan bool)

	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				do(event)
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	flag.Parse()
	root := flag.Arg(0)
	err = filepath.Walk(root, walk)
	fmt.Printf("filepath.Walk() returned %v\n", err)

	<-done
}
