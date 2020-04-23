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
	var fromID, toID uint32
	var indexOfToID int
	for scanner.Scan() {
		line = scanner.Text()
		lineSplit = strings.Split(line, ",")
		fromID = MustParseBase36(lineSplit[0])
		for _, toIDString = range lineSplit[1:] {
			toID = MustParseBase36(toIDString)
			indexOfToID = pageIDs.SearchMustFind(toID)
			backlinks[indexOfToID] = append(backlinks[indexOfToID], fromID)
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
	for i, id := range pageIDs {
		outWriter.WriteString(strconv.FormatInt(int64(id), 36))
		for _, fromID := range backlinks[i] {
			outWriter.WriteString(",")
			outWriter.WriteString(strconv.FormatInt(int64(fromID), 36))
		}
		outWriter.WriteString("\n")
	}
	return nil
}
