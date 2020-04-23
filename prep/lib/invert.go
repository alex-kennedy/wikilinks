package lib

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

//SaveBacklinks inverts the pivoted pagelinks file and produces a file of the
//same form, instead representing backlinks.
func SaveBacklinks(pagelinksPivotedName, pageIDsName, backlinksName string) error {
	pageIDs, err := LoadPageIDs(pageIDsName)
	if err != nil {
		return err
	}

	pagelinksPivotedFile, err := os.Open(pagelinksPivotedName)
	if err != nil {
		return err
	}
	defer pagelinksPivotedFile.Close()
	scanner := bufio.NewScanner(pagelinksPivotedFile)

	bar := NewProgressBarFileSize(pagelinksPivotedFile)

	backlinks := make([][]uint32, len(pageIDs))

	var line, toIDString string
	var lineSplit []string
	var fromIndex, toIndex uint32
	for scanner.Scan() {
		line = scanner.Text()
		lineSplit = strings.Split(line, ",")
		fromIndex = MustParseBase36(lineSplit[0])
		for _, toIDString = range lineSplit[1:] {
			toIndex = MustParseBase36(toIDString)
			backlinks[toIndex] = append(backlinks[toIndex], fromIndex)
		}
		bar.Add(len(line) + 1)
	}
	bar.Finish()

	log.Println("Writing to disk...")
	backlinksFile, err := os.Create(backlinksName)
	if err != nil {
		return err
	}
	defer backlinksFile.Close()
	outWriter := bufio.NewWriter(backlinksFile)
	defer outWriter.Flush()
	for toIndex := range pageIDs {
		outWriter.WriteString(strconv.FormatInt(int64(toIndex), 36))
		for _, fromIndex := range backlinks[toIndex] {
			outWriter.WriteString(",")
			outWriter.WriteString(strconv.FormatInt(int64(fromIndex), 36))
		}
		outWriter.WriteString("\n")
	}
	return nil
}
