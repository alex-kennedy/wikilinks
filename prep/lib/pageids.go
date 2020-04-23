//Functions related to loading, parsing, writing, and searching page IDs.

package lib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

//SavePageIDs gets a list of all the possible page IDs from page_merged (base 10
//IDs), sorts them, and writes them to a file (base 36 IDs). This file
//represents a complete list of page IDs that exist in data dump. This file
//stores real Wikipedia page IDs. A value which indexes into this file is
//referred to as a 'page index' in this library.
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

//PageIDs is a searchable array of the real Wikipedia page IDs. It is used to
//acquire page indexes and is the canonical mapping of index to ID for the data
//dump in use. Note that all page IDs which make it to past the resolution steps
//must appear in this list. The set up of the WikiMedia SQL server means each
//page ID will definitely fit in a uint32.
type PageIDs []uint32

//SearchMustFind finds the index of a page which must definitely be a legal page
//ID.
func (p PageIDs) SearchMustFind(x uint32) int {
	i := sort.Search(len(p), func(i int) bool { return p[i] >= x })
	if i < len(p) && p[i] == x {
		return i
	}
	log.Fatalf("Page ID %d could not be found in page.", x)
	return -1
}

//LoadPageIDs loads the canonical, sorted page IDs file.
func LoadPageIDs(in string) (PageIDs, error) {
	fIn, err := os.Open(in)
	if err != nil {
		return nil, err
	}
	defer fIn.Close()
	scanner := bufio.NewScanner(fIn)

	bar := NewProgressBarFileSize(fIn)
	defer bar.Finish()

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

//MustParseBase36 converts an index or page ID which must be a valid base-36
//string.
func MustParseBase36(idString string) uint32 {
	id, err := strconv.ParseUint(idString, 36, 32)
	if err != nil {
		log.Fatalf("%s could not be passed to uint32", idString)
	}
	return uint32(id)
}
