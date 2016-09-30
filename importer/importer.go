package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/ogier/pflag"
	"gopkg.in/olivere/elastic.v3"
)

var (
	tagTitle       = regexp.MustCompile("(?i)TAG:TITLE=(.*)")
	tagArtist      = regexp.MustCompile("(?i)TAG:ARTIST=(.*)")
	tagAlbum       = regexp.MustCompile("(?i)TAG:ALBUM=(.*)")
	tagTrack       = regexp.MustCompile("(?i)TAG:TRACK=(.*)")
	tagGenre       = regexp.MustCompile("(?i)TAG:GENRE=(.*)")
	tagDate        = regexp.MustCompile("(?i)TAG:DATE=(.*)")
	tagAlbumArtist = regexp.MustCompile("(?i)TAG:ALBUM[_| ]ARTIST=(.*)")
	tagDuration    = regexp.MustCompile("(?i)DURATION=(.*)")
	tagFileName    = regexp.MustCompile("(?i)FILENAME=(.*)")
	tagSize        = regexp.MustCompile("(?i)SIZE=(.*)")
	tagBitRate     = regexp.MustCompile("(?i)BIT[_| ]RATE=(.*)")

	ffmpegBin  string
	ffprobeBin string
	coversPath string

	ds = string(filepath.Separator)

	// https://en.wikipedia.org/wiki/Audio_file_format
	extensionsList = []string{
		//".3gp",
		//".aa",
		//".aac",
		//".aax",
		//".act",
		//".aiff",
		//".amr",
		//".ape",
		//".au",
		//".awb",
		//".dct",
		//".dss",
		//".dvf",
		//".flac",
		//".gsm",
		//".iklax",
		//".ivs",
		//".m4a",
		//".m4b",
		//".m4p",
		//".mmf",
		".mp3",
		//".mpc",
		//".msv",
		//".ogg",
		//".oga",
		//".opus",
		//".ra",
		//".rm",
		//".raw",
		//".sln",
		//".tta",
		//".vox",
		//".wav",
		//".wma",
		//".wv",
		//".webm",
	}

	client *elastic.Client
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	path := pflag.String("path", "", "path to import")
	pflag.Parse()
	//TODO: check path exists

	var err error
	client, err = elastic.NewClient()
	assertNil(err)

	// Delete an index.
	//deleteIndex, err := client.DeleteIndex("song").Do()
	//assertNil(err)
	//if !deleteIndex.Acknowledged {
	//	// Not acknowledged
	//}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("song").Do()
	assertNil(err)
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("song").Do()
		assertNil(err)
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	//wd, err := os.Getwd()
	//assertNil(err)

	//ffmpegBin = wd + "ffmpeg/ffmpeg"
	ffmpegBin = "ffmpeg"
	//ffprobeBin = wd + "/ffmpeg/ffprobe"
	ffprobeBin = "ffprobe"
	//coversPath = wd + "/.thumbs"

	//assertNil(os.MkdirAll(coversPath, 0755))

	logrus.Debug(*path)
	assertNil(filepath.Walk(*path, walkFn))
}

func walkFn(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return nil
	}

	if fi.IsDir() {
		return nil
	}

	var extension = filepath.Ext(fi.Name())

	if !isValidExtension(extension) {
		return nil
	}

	//if _, err = models.GetSongBySongPath(path); err == nil {
	//Song Already Exists
	// TODO: must sync
	//return nil
	//}

	logrus.Info(path)

	metaString, err := fetchMetaData(path)
	if err != nil {
		logrus.Warn(err)
		return nil
	}

	song := map[string]interface{}{}

	if t := tagTitle.FindStringSubmatch(metaString); len(t) > 1 {
		song["title"] = t[1]
	}

	if t := tagArtist.FindStringSubmatch(metaString); len(t) > 1 {
		song["artist"] = t[1]
	}

	if t := tagAlbum.FindStringSubmatch(metaString); len(t) > 1 {
		song["album"] = t[1]
	}

	if t := tagAlbumArtist.FindStringSubmatch(metaString); len(t) > 1 {
		song["album_artist"] = t[1]
	}

	if t := tagDate.FindStringSubmatch(metaString); len(t) > 1 {
		song["date"] = t[1]
	}

	if t := tagGenre.FindStringSubmatch(metaString); len(t) > 1 {
		song["genre"] = t[1]
	}

	if t := tagFileName.FindStringSubmatch(metaString); len(t) > 1 {
		song["file_name"] = t[1]
	}

	if t := tagSize.FindStringSubmatch(metaString); len(t) > 1 {
		song["size"] = t[1]
	}

	if t := tagBitRate.FindStringSubmatch(metaString); len(t) > 1 {
		song["bit_rate"] = t[1]
	}

	if t := tagTrack.FindStringSubmatch(metaString); len(t) > 1 {
		song["track"] = t[1]
	}

	if t := tagDuration.FindStringSubmatch(metaString); len(t) > 1 {
		//duration, err := ReadableTime(strings.TrimSpace(t[1]))
		//if err != nil {
		//	logrus.Warn(err)
		//} else {
		//	song.Duration = duration
		//}
		song["duration"] = t[1]
	}

	//song["path"] = path

	//err = fetchCover(path)
	//if err != nil {
	//	logrus.Warn(err)
	//} else {
	//	song["cover"] = path + ".jpg"
	//}

	//if cover.Valid {
	//	if err = createImageSize(thumbnailsPath, 300, 300); err != nil {
	//		return err
	//	}
	//
	//	if err = createImageSize(thumbnailsPath, 500, 500); err != nil {
	//		return err
	//	}
	//
	//	if err = createImageSize(thumbnailsPath, 900, 900); err != nil {
	//		return err
	//	}
	//}
	hasher := sha256.New()
	f, err := os.Open(path)
	if err != nil {
		logrus.Warn(err)
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatal(err)
	}
	hashString := fmt.Sprintf("%x", hasher.Sum(nil))

	_, err = client.Index().
		Index("songs").
		Type("song").
		Id(hashString).
		BodyJson(song).
		Refresh(true).
		Do()
	assertNil(err)

	return nil
}

func isValidExtension(ext string) bool {
	for _, v := range extensionsList {
		if ext == v {
			return true
		}
	}

	return false
}

func fetchMetaData(input string) (string, error) {
	cmd := exec.Command(ffprobeBin, "-show_format", input)
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		logrus.Warn(err)
		logrus.Warn(stdErr.String())
		logrus.Warn(stdOut.String())
		return "", err
	}
	return stdOut.String(), nil
}

func fetchCover(input string) error {
	cmd := exec.Command(ffmpegBin, "-i", input, "-an", "-vcodec", "copy", input+".jpg")
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	return cmd.Run()
}

func createImageSize(input string, w, h int) error {
	var extension = filepath.Ext(input)
	var name = input[0 : len(input)-len(extension)]
	output := fmt.Sprintf("%s_%dx%d%s", name, w, h, extension)

	cmd := exec.Command(ffmpegBin, "-i", input, "-vf", fmt.Sprintf("scale=%d:%d", w, h), output)
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		logrus.Warn(err)
		logrus.Warn(stdErr.String())
		return err
	}
	logrus.Info(stdOut.String())

	return nil
}

// assertNil panic if the test is not nil
func assertNil(test interface{}, params ...interface{}) {
	if test != nil {
		f := logrus.Fields{}
		for i := range params {
			f[fmt.Sprintf("param%d", i)] = params[i]
		}

		if e, ok := test.(error); ok {
			logrus.WithFields(f).Panic(e)
		}
		logrus.WithFields(f).Panic("must be nil, but its not")
	}
}

// ReadableTime make time in readable format
func ReadableTime(s string) (string, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%02d:%02d", int64(f/60), int64(math.Mod(f, 60))), nil
}
