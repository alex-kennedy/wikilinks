package lib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb"
)

//Searcher represents a set of values that can be looked up by key. Here, it is
//implemented as an external binary search, as a map, and as a few specialised
//in memory binary searches.
type Searcher interface {
	Search(string) (string, error)
}

//NotFoundError is returned when the key is not found.
type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return fmt.Sprint("search key not found")
}

//
//Binary search.
//

//BinarySearcher holds the information needed to binary search a sorted file.
type BinarySearcher struct {
	file           *os.File
	index          []int64
	buffer         []byte
	keyValFunction KeyValFunction
}

//createReadBuffer creates a buffer large enough to hold any one line.
func createReadBuffer(index []int64) []byte {
	capNeeded := int64(0)
	for i := range index[:len(index)-1] {
		if index[i+1]-index[i] > capNeeded {
			capNeeded = index[i+1] - index[i]
		}
	}
	return make([]byte, capNeeded)
}

func (b *BinarySearcher) getLine(i int) (string, string) {
	b.file.Seek(b.index[i], 0)
	buffer := b.buffer[:b.index[i+1]-b.index[i]-1]
	n, _ := b.file.Read(buffer)
	if n == 0 {
		return "", ""
	}
	line := string(buffer)
	return b.keyValFunction(line)
}

func (b *BinarySearcher) binarySearch(key *string, lo, hi int) (string, error) {
	if lo > hi {
		return "", &NotFoundError{}
	}
	middle := (lo + hi) / 2

	k, v := b.getLine(middle)
	if k == *key {
		return v, nil
	}

	if *key > k {
		return b.binarySearch(key, middle+1, hi)
	}
	return b.binarySearch(key, lo, middle-1)
}

//Search runs a binary search for the given key.
func (b *BinarySearcher) Search(key string) (string, error) {
	//The length of the file is one less than the index, and the last index is
	//one less again.
	return b.binarySearch(&key, 0, len(b.index)-2)
}

//NewBinarySearcher creates an object for binary searching a file. Requires a
//sorted, indexed file.
func NewBinarySearcher(name, indexName string, keyVal KeyValFunction) (*BinarySearcher, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	index, err := loadIndex(indexName)
	if err != nil {
		return nil, err
	}
	buffer := createReadBuffer(index)
	return &BinarySearcher{f, index, buffer, keyVal}, nil
}

//
//Map search.
//

//MapSearcher loads the complete file into a hash map for fast, in-memory
//searching.
type MapSearcher struct {
	keyValMap map[string]string
}

//Search finds the value associated with key or returns a NotFoundError.
func (m *MapSearcher) Search(key string) (string, error) {
	value, ok := m.keyValMap[key]
	if ok {
		return value, nil
	}
	return "", &NotFoundError{}
}

//NewMapSearcher creates a map of all the key, value pairs in the given file.
func NewMapSearcher(fileName string, keyVal KeyValFunction) (*MapSearcher, error) {
	log.Printf("Mapping %s...", fileName)
	bar := pb.StartNew(-1)
	defer bar.Finish()

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	keyValMap := make(map[string]string)
	for scanner.Scan() {
		key, value := keyVal(scanner.Text())
		keyValMap[key] = value
		bar.Add(1)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return &MapSearcher{keyValMap}, nil
}

//
//Array search.
//

//StringToIntArraySearcher loads a file of line format "string,int" and allows
//searching it by title.
type StringToIntArraySearcher struct {
	keys   []string
	values []int
}

//NewStringToIntArraySearcher creates the container for searching the given
//file. The file is loaded into memory as two arrays.
func NewStringToIntArraySearcher(fileName string) (*StringToIntArraySearcher,
	error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var keys []string
	var values []int

	bar := pb.StartNew(-1)
	defer bar.Finish()

	scanner := bufio.NewScanner(file)
	var line, key, value string
	var valueInt int
	for scanner.Scan() {
		line = scanner.Text()
		key, value = KeyValLastComma(line)
		valueInt, err = strconv.Atoi(value)
		if err != nil {
			continue
		}
		keys = append(keys, string([]byte(key)))
		values = append(values, valueInt)
		bar.Add(1)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	runtime.GC()
	return &StringToIntArraySearcher{keys, values}, nil
}

//Search searches the loaded file by key. Returns -1 if the search fails.
func (s *StringToIntArraySearcher) Search(key string) int {
	i := sort.SearchStrings(s.keys, key)
	if i >= len(s.keys) || s.keys[i] != key {
		return -1
	}
	return s.values[i]
}

//PagelinksPivoted is for obtaining the list of edges originating at each page.
type PagelinksPivoted interface {
	GetLinks(int) []uint32
	GetNumberOfPages() int
}

//PagelinksPivotedInMemory represents a loaded pagelinks_pivoted file for
//searching.
type PagelinksPivotedInMemory struct {
	Destinations  [][]uint32
	NumberOfPages int
}

//GetLinks for the in memory case is essentially an index.
func (p *PagelinksPivotedInMemory) GetLinks(index int) []uint32 {
	return p.Destinations[index]
}

//GetNumberOfPages gets the total number of pages.
func (p *PagelinksPivotedInMemory) GetNumberOfPages() int {
	return p.NumberOfPages
}

//LinkExists returns whether the link from -> to exists.
func (p *PagelinksPivotedInMemory) LinkExists(from, to uint32) bool {
	outs := p.Destinations[from]
	i := sort.Search(len(outs), func(i int) bool { return outs[i] >= to })
	return i < len(outs) && outs[i] == to
}

//NewPagelinksPivotedInMemory indexes a pivoted pagelink file for easy
//searching. Note that this loads the pagelinks file into memory.
func NewPagelinksPivotedInMemory(pagelinksName, pageIDsName string) (*PagelinksPivotedInMemory, error) {
	pagelinksFile, err := os.Open(pagelinksName)
	if err != nil {
		return nil, err
	}
	defer pagelinksFile.Close()
	scanner := bufio.NewScanner(pagelinksFile)
	scanner.Buffer(make([]byte, 16*1024*1024), 16*1024*1024) // 16MB

	bar := NewProgressBarFileSize(pagelinksFile)
	defer bar.Finish()

	pageIDs, err := LoadPageIDs(pageIDsName)
	if err != nil {
		return nil, err
	}

	pagesWithAtLeastOneLink := 0

	destinations := make([][]uint32, len(pageIDs))
	var line string
	var lineSplit []string
	var source uint32
	var lineDestinations []uint32
	for scanner.Scan() {
		line = scanner.Text()
		bar.Add(len(line) + 1) //File is UTF-8 ASCII, newlines are dropped
		lineSplit = strings.Split(line, ",")
		source = MustParseBase36(lineSplit[0])
		lineDestinations = make([]uint32, len(lineSplit)-1)
		for i, pageString := range lineSplit[1:] {
			lineDestinations[i] = MustParseBase36(pageString)
		}
		if len(lineDestinations) > 0 {
			pagesWithAtLeastOneLink++
		}
		destinations[source] = lineDestinations
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return &PagelinksPivotedInMemory{destinations, len(destinations)}, nil
}
