package lib

import (
	"bufio"
	"os"
	"sort"
	"strings"

	"github.com/cheggaaa/pb"
)

//PivotFile takes a sorted file of key,value and puts all values from the same
//key on the same line. For example, k1,v1; k1,v2 becomes k1,v1,v2.
func PivotFile(inPath, outPath string, bytesPerBuffer int) error {
	inFile, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scannerBuffer := make([]byte, bytesPerBuffer)
	scanner.Buffer(scannerBuffer, bytesPerBuffer)

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	outWriter := bufio.NewWriterSize(outFile, bytesPerBuffer)
	defer outWriter.Flush()

	bar := NewProgressBarFileSize(inFile)
	defer bar.Finish()

	doPivot(scanner, outWriter, bar)
	return nil
}

func doPivot(in *bufio.Scanner, out *bufio.Writer, bar *pb.ProgressBar) {
	in.Scan()
	line := in.Text()
	k, v := KeyValFirstComma(line)
	mainKey := k
	values := []string{v}

	for in.Scan() {
		line := in.Text()
		k, v := KeyValFirstComma(line)
		if k != mainKey {
			sort.Strings(values)
			out.WriteString(mainKey + "," + strings.Join(values, ",") + "\n")
			mainKey = k
			values = values[:0]
		}
		values = append(values, v)
		bar.Add(len(line) + 1)
	}
}
