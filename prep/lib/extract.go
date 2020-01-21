package lib

import (
	"bufio"
	"encoding/csv"
	"os"
	"strings"

	"github.com/cheggaaa/pb"
)

func toCsvEscaping(record string) string {
	builder := strings.Builder{}
	escaped := false
	for _, r := range record {
		if escaped {
			if r == '"' {
				builder.WriteString("\"\"")
			} else {
				builder.WriteRune(r)
			}
			escaped = false
		} else if r == '\\' {
			escaped = true
		} else if r == '\'' {
			builder.WriteRune('"')
		} else {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func processLine(line string, fieldsPerRecord int) ([][]string, error) {
	subString := line[strings.IndexRune(line, '(')+1 : len(line)-2]
	entries := strings.Split(subString, "),(")
	for i, entry := range entries {
		entries[i] = toCsvEscaping(entry)
	}
	csvDocument := strings.Join(entries, "\n")

	csvReader := csv.NewReader(strings.NewReader(csvDocument))
	csvReader.FieldsPerRecord = fieldsPerRecord

	v, err := csvReader.ReadAll()
	return v, err
}

//ExtractTable extracts the given column indices from the entries in inPath.
func ExtractTable(inPath, outPath string, indices []int, fieldsPerRecord int) error {
	fIn, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer fIn.Close()

	scanner := bufio.NewScanner(fIn)
	const maxCapacity = 10 * 1024 * 1024 // 10MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	fOut, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer fOut.Close()
	writer := csv.NewWriter(fOut)
	output := make([]string, len(indices))

	pb := pb.StartNew(-1)

	dataRemaining := true
	for dataRemaining {
		dataRemaining = scanner.Scan()
		if scanner.Err() != nil {
			return scanner.Err()
		}
		pb.Add(1)
		line := scanner.Text()
		if !strings.HasPrefix(line, "INSERT INTO") {
			continue
		}
		records, err := processLine(line, fieldsPerRecord)
		if err != nil {
			return err
		}
		for _, record := range records {
			for _, i := range indices {
				output[i] = record[i]
			}
			err := writer.Write(output)
			if err != nil {
				return err
			}
		}
	}
	pb.Finish()
	return nil
}
