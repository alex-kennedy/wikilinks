package lib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb"
)

//SavePageIDs gets a list of all the possible page IDs from page_merged (base 10
//IDs), sorts them, and writes them to a file (base 36 IDs).
func SavePageIDs(in, out string) error {
	fIn, err := os.Open(in)
	if err != nil {
		return err
	}
	defer fIn.Close()
	scanner := bufio.NewScanner(fIn)

	ids := []int{}
	var line, idString string
	var id int
	for scanner.Scan() {
		line = scanner.Text()
		_, idString = KeyValLastComma(line)
		id, err = strconv.Atoi(idString)
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}

	sort.Ints(ids)

	fOut, err := os.Create(out)
	if err != nil {
		return err
	}
	defer fOut.Close()
	outWriter := bufio.NewWriter(fOut)
	for _, id := range ids {
		outWriter.WriteString(fmt.Sprintf("%s\n", strconv.FormatInt(int64(id), 36)))
	}
	defer outWriter.Flush()
	return nil
}

func mustParsePageIDBase36(idString string) uint32 {
	id, err := strconv.ParseUint(idString, 36, 32)
	if err != nil {
		log.Fatalf("%s could not be passed to uint32", idString)
	}
	return uint32(id)
}

type PageIDs []uint32

func (p PageIDs) SearchMustFind(x uint32) int {
	i := sort.Search(len(p), func(i int) bool { return p[i] >= x })
	if i < len(p) && p[i] == x {
		return i
	}
	log.Fatalf("Page ID %d could not be found in page.", x)
	return -1
}

func LoadPageIDs(in string) (PageIDs, error) {
	fIn, err := os.Open(in)
	if err != nil {
		return nil, err
	}
	defer fIn.Close()
	scanner := bufio.NewScanner(fIn)

	ids := []uint32{}
	var idString string
	for scanner.Scan() {
		idString = scanner.Text()
		id, err := strconv.ParseUint(idString, 36, 32)
		if err != nil {
			return nil, err
		}
		ids = append(ids, uint32(id))
	}
	return ids, nil
}

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

	fileInfo, err := pagelinksPivotedFile.Stat()
	if err != nil {
		return err
	}
	bar := pb.Start64(fileInfo.Size())
	bar.Set(pb.Bytes, true)

	backlinks := make([][]uint32, len(pageIDs))

	var line, toIDString string
	var lineSplit []string
	var fromID, toID uint32
	var indexOfToID int
	for scanner.Scan() {
		line = scanner.Text()
		lineSplit = strings.Split(line, ",")
		fromID = mustParsePageIDBase36(lineSplit[0])
		for _, toIDString = range lineSplit[1:] {
			toID = mustParsePageIDBase36(toIDString)
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
	for i, id := range pageIDs {
		outWriter.WriteString(strconv.FormatInt(int64(id), 36))
		for _, fromID := range backlinks[i] {
			outWriter.WriteString(",")
			outWriter.WriteString(strconv.FormatInt(int64(fromID), 36))
		}
		outWriter.WriteString("\n")
	}
	outWriter.Flush()
	return nil
}
