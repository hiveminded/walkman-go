package id3

import (
	id3 "github.com/mikkyang/id3-go"
	id3v2 "github.com/mikkyang/id3-go/v2"
)

func GetUniqueFileIdentifier(path string) (string, error) {
	mp3File, err := id3.Open(path)
	defer mp3File.Close()
	if err != nil {
		return "", err
	}

	return mp3File.Frame("UFID").String(), nil
}

func SetUniqueFileIdentifier(path, ufid string) error {
	mp3File, err := id3.Open(path)
	defer mp3File.Close()
	if err != nil {
		return err
	}

	ft := id3v2.V23FrameTypeMap["UFID"]

	textFrame := id3v2.NewTextFrame(ft, ufid)
	mp3File.AddFrames(textFrame)

	return nil
}

type Track struct {
	UFID     string   `json:"ufid"`
	Title    string   `json:"title"`
	Artist   string   `json:"artist"`
	Album    string   `json:"album"`
	Year     string   `json:"year"`
	Genre    string   `json:"genre"`
	Comments []string `json:"comments"`
}

func (obj *Track) GetID() string {
	return obj.UFID
}

func (obj *Track) GetType() string {
	return "track"
}

func GetFileInfo(path string) (*Track, error) {
	mp3File, err := id3.Open(path)
	defer mp3File.Close()
	if err != nil {
		return nil, err
	}

	track := &Track{}
	track.UFID = mp3File.Frame("UFID").String()
	track.Title = mp3File.Title()
	track.Artist = mp3File.Artist()
	track.Album = mp3File.Album()
	track.Year = mp3File.Year()
	track.Genre = mp3File.Genre()
	track.Comments = mp3File.Comments()

	return track, nil
}
