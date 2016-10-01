package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func watch() {
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
}

func do(event fsnotify.Event) {
	log.Println("event:", event)
	if event.Op&fsnotify.Write == fsnotify.Write {
		log.Println("modified file:", event.Name)
	}

	//Switch event {
	//	deleteFolder:
	//	Delete all items from list and index
	//	Deletefile:
	//	Delete item form list and index
	//	Renamefolder:
	//	Update all items in list and index
	//	Renamefile:
	//	Update item in list and index
	//	Addfolder:
	//	Add all items in list and index
	//	Addfile:
	//	Add item in list and index
	//	Changefile
	//	Update item in list and index
	//}
}
