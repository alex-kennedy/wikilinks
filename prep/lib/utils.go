package lib

import (
	"os"

	"github.com/cheggaaa/pb"
)

//CleanupFile will delete the file if it exists, returning any errors.
func CleanupFile(file string) error {
	err := os.Remove(file)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

//CheckExists returns true if the file exists, false otherwise or on an error.
func CheckExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

//writeCounter is a helper struct for printing progress during the file
//download.
type writeCounter struct {
	pb *pb.ProgressBar
}

func (wc writeCounter) New(count int64) writeCounter {
	wc.pb = pb.StartNew(int(count))
	wc.pb.Set(pb.Bytes, true)
	return wc
}

func (wc writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.pb.Add(n)
	return n, nil
}
