// manager is a service that will update elastic database with music directory file change
/*
Workflow

`
StartUp {
	WalkPath {
		if IsDir then add to watch
		if IsFile then {
			Check or create in file list
			If is new or hash changed then insert in index
		}
	}

	getNotCheckedFiles().then{
		Remove from index and file list
	}
}

Watch {
	Switch event {
		DeleteFolder:
			Delete all items from list and index
		DeleteFile:
			Delete item form list and index
		RenameFolder:
			Update all items in list and index
		RenameFile:
			Update item in list and index
		AddFolder:
			Add all items in list and index
		AddFile:
			Add item in list and index
		ChangeFile:
			Update item in list and index
	}
}
`
*/

package main

import (
	"flag"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

var (
	watcher *fsnotify.Watcher
	root    string = os.Getenv("HOME")
)

func main() {
	var err error

	// new fsnotify watcher
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	//set root path
	flag.Parse()
	if len(flag.Arg(0)) > 0 {
		root = flag.Arg(0)
	}

	// do start up check
	startUp()

	// watch directories
	watch()

	// wait for interrupt signal
	<-make(chan bool)
}
