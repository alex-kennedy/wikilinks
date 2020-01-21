package lib

import (
	"compress/gzip"
	"io"
	"os"
)

//UnzipGzFile extracts a .gz file to outPath.
func UnzipGzFile(inPath, outPath string) error {
	fIn, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer fIn.Close()

	gzReader, err := gzip.NewReader(fIn)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	fOut, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer fOut.Close()

	counter := writeCounter{}.New(int64(-1))
	_, err = io.Copy(fOut, io.TeeReader(gzReader, counter))
	if err != nil {
		return err
	}
	counter.pb.Finish()
	return nil
}
