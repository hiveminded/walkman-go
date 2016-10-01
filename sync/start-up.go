package main

import (
	"github.com/lujem/walkman-go/id3"
	"github.com/rs/xid"
	"log"
	"os"
	"path/filepath"
)

func startUp() {
	err := filepath.Walk(root, walk)
	if err != nil {
		log.Fatal(err)
	}

	//getNotcheckedfiles().then{
	//	Remove from index and filelist
	//}
}

func walk(path string, f os.FileInfo, err error) error {

	// if is directory, only add to watcher
	if f.IsDir() {
		return watcher.Add(path)
	}

	// ok, this is a file
	_, err = fetchOrInitUniqueFileIdentifier(path)

	track, err := id3.GetFileInfo(path)
	if err != nil {
		return err
	}

	//Check or create in filelist
	//If is new or hash changed then insert in index

	return nil
}

// registerFile will assign an id to file and store it in file id3 tags
// http://id3.org/id3v2.3.0#Unique_file_identifier
func fetchOrInitUniqueFileIdentifier(path string) (string, error) {
	// check ufid exists
	ufid, err := id3.GetUniqueFileIdentifier(path)
	if err != nil {
		return "", err
	}

	if ufid == "" {
		// generate new id
		guid := xid.New()
		ufid = guid.String()

		// assign new id to file
		err = id3.SetUniqueFileIdentifier(path, ufid)
		if err != nil {
			return "", err
		}
	}

	return ufid, nil
}
