package lib

import (
	"fmt"
	"os"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

//CleanupFile will delete the file if it exists, returning any errors.
func CleanupFile(file string) error {
	err := os.Remove(file)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

//CheckExists returns true if the file exists, false otherwise or on an error.
func CheckExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

//writeCounter is a helper struct for printing progress during the file
//download.
type writeCounter struct {
	pb *pb.ProgressBar
}

func (wc writeCounter) New(count int64) writeCounter {
	wc.pb = pb.StartNew(int(count))
	wc.pb.Set(pb.Bytes, true)
	return wc
}

func (wc writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.pb.Add(n)
	return n, nil
}

//KeyValFunction returns the key and value from a file line for sorting and
//searching.
type KeyValFunction func(string) (string, string)

//KeyValLastComma splits a line into key and value based on the last comma,
//suitable when the ID is last.
func KeyValLastComma(s string) (string, string) {
	commaIndex := strings.LastIndex(s, ",")
	if commaIndex == -1 {
		return s, ""
	}
	return string([]byte(s[:commaIndex])), string([]byte(s[commaIndex+1:]))
}

//KeyValFirstComma splits a line into key and value based on the first comma,
//suitable when the ID is first.
func KeyValFirstComma(s string) (string, string) {
	commaIndex := strings.Index(s, ",")
	if commaIndex == -1 {
		return s, ""
	}
	return string([]byte(s[:commaIndex])), string([]byte(s[commaIndex+1:]))
}

//NewProgressBarFileSize creates a new progress bar based on the size of the
//file. When reading from the file, add the number of bytes read. Remember to
//Finish() the progress bar. This function panics if it can't read the size of
//the file.
func NewProgressBarFileSize(file *os.File) *pb.ProgressBar {
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	bar := pb.Start64(fileInfo.Size())
	bar.Set(pb.Bytes, true)
	return bar
}

//GetWikiNameFromURL determines the wiki name ("enwiki" or "simplewiki") from
//the dump URL. It should be the last component of the name.
func GetWikiNameFromURL(url string) (string, error) {
	urlBits := strings.Split(strings.TrimRight(url, "/"), "/")
	wikiName := urlBits[len(urlBits)-1]
	if wikiName == "enwiki" || wikiName == "simplewiki" {
		return wikiName, nil
	}
	return "", fmt.Errorf("invalid wiki name: %s", wikiName)
}
