package lib

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/cheggaaa/pb"
)

//IndexFile produces a file of int64s with the byte locations of the beginning
//of each line. This function is very inefficient, but not enough to compel me
//to make improvements.
func IndexFile(in, out string) error {
	bar := pb.StartNew(-1)
	defer bar.Finish()

	pos := []int64{0}
	nextByte := []byte{0}
	newLine := []byte("\n")

	fIn, err := os.Open(in)
	if err != nil {
		return err
	}
	defer fIn.Close()

	for {
		n, err := fIn.Read(nextByte)
		i, err := fIn.Seek(0, os.SEEK_CUR)
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}
		//Though this is a UTF8 file with variable length characters, it is so that
		//this byte pair will never occur except to represent a new line. Nifty.
		if nextByte[0] != newLine[0] {
			continue
		}
		pos = append(pos, i)
		bar.Add(1)
	}

	fOut, err := os.Create(out)
	if err != nil {
		return err
	}
	defer fOut.Close()
	for _, n := range pos {
		fOut.WriteString(fmt.Sprintln(n))
	}

	return nil
}

func loadIndex(name string) ([]int64, error) {
	fIndex, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fIndex.Close()

	scanner := bufio.NewScanner(fIndex)

	index := make([]int64, 0)
	for scanner.Scan() {
		n, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		index = append(index, n)
	}
	return index, nil
}

//NotFoundError is returned when the key is not found in the file.
type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return fmt.Sprint("search key not found")
}

//getReadBuffer creates a buffer large enough to hold any one line.
func getReadBuffer(index []int64) []byte {
	capNeeded := int64(0)
	for i := range index[:len(index)-1] {
		if index[i+1]-index[i] > capNeeded {
			capNeeded = index[i+1] - index[i]
		}
	}
	return make([]byte, capNeeded)
}

//BinarySearcher holds the information needed to binary search a sorted file.
type BinarySearcher struct {
	file           *os.File
	index          []int64
	buffer         []byte
	keyValFunction KeyValFunction
}

//Search runs a binary search for the given key.
func (b *BinarySearcher) Search(key string) (int, string, error) {
	//The length of the file is one less than the index, and the last index is
	//one less again.
	return b.binarySearch(&key, 0, len(b.index)-2)
}

func (b *BinarySearcher) binarySearch(key *string, lo, hi int) (int, string, error) {
	if lo > hi {
		return -1, "", &NotFoundError{}
	}
	middle := (lo + hi) / 2

	k, v := b.getLine(middle)
	if k == *key {
		return middle, v, nil
	}

	if *key > k {
		return b.binarySearch(key, middle+1, hi)
	}
	return b.binarySearch(key, lo, middle-1)
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

//NewBinarySearcher creates an object for binary searching a file.
func NewBinarySearcher(name, indexName string, keyVal KeyValFunction) (*BinarySearcher, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	index, err := loadIndex(indexName)
	if err != nil {
		return nil, err
	}
	buffer := getReadBuffer(index)
	return &BinarySearcher{f, index, buffer, keyVal}, nil
}
