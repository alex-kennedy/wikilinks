package lib

import (
	"fmt"
	"github.com/cheggaaa/pb"
	"os"
)

//IndexFile produces a file of int64s with the byte locations of the beginning
//of each line.
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
