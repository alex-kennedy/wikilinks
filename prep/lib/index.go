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

//loadIndex loads an index file to an array of byte-marks.
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
