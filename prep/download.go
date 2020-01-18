package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Jeffail/gabs/v2"

	"github.com/dustin/go-humanize"
)

//verifyWikiFile accesses the MD5 of file as reported by the mirror and
//compares it to the MD5 hash of the downloaded file.
func verifyWikiFile(fileName, siteURL, outPath string) error {
	response, err := http.Get(siteURL + "dumpstatus.json")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	dumpStatus, err := gabs.ParseJSON(data)

	fileNamePieces := strings.Split(strings.Split(fileName, ".")[0], "-")
	tableName := fileNamePieces[len(fileNamePieces)-1] + "table"
	correctHash := dumpStatus.Search("jobs", tableName, "files", fileName, "md5").String()
	if len(correctHash) != 34 {
		return errors.New("Failed to get valid MD5 from JSON, got " + correctHash)
	}
	correctHash = correctHash[1 : len(correctHash)-1] // Accounts for quotes

	file, err := os.Open(outPath)
	if err != nil {
		return err
	}
	defer file.Close()
	myFileHash := md5.New()
	if _, err := io.Copy(myFileHash, file); err != nil {
		return err
	}
	myFileHashString := hex.EncodeToString(myFileHash.Sum(nil))

	if correctHash != myFileHashString {
		return errors.New("Hashes not equal")
	}

	return nil
}

//writeCounter is a helper struct for printing progress during the file
//download.
type writeCounter struct {
	Total uint64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc writeCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 50))
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

//DownloadWikiFile downloads and verifies fileName from the wiki dump at
//siteUrl. Saves it to outPath.
func DownloadWikiFile(fileName, siteURL, outPath string) error {
	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	response, err := http.Get(siteURL + fileName)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", response.Status)
	}

	counter := &writeCounter{}
	_, err = io.Copy(out, io.TeeReader(response.Body, counter))
	if err != nil {
		return err
	}

	err = verifyWikiFile(fileName, siteURL, outPath)
	if err != nil {
		return err
	}
	fmt.Println()
	return nil
}
