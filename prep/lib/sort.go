package lib

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/cheggaaa/pb/v3"
)

//chunk represents one sorted file
type chunk struct {
	head           string
	key            string
	reader         *bufio.Scanner
	alive          bool
	keyValFunction KeyValFunction
}

func (c *chunk) pop() {
	c.alive = c.reader.Scan()
	c.head = c.reader.Text()
	k, _ := c.keyValFunction(c.head)
	c.key = k
}

func (c *chunk) lessThan(d *chunk) bool {
	return c.key < d.key
}

//parent calculates the parent node of the ith node.
func parent(i int) int {
	if i == 0 {
		return -1
	}
	return int((i - 1) / 2)
}

//child calculates the left child of the ith node.
func child(i int) int {
	return (2 * i) + 1
}

//fheap implements the min-heap of sorted file chunks.
type fheap struct {
	n    int
	heap []*chunk
	out  *bufio.Writer
}

func (h *fheap) insert(c *chunk) {
	h.n++
	h.heap = append(h.heap, c)
	h.bubbleUp(h.n - 1)
}

func (h *fheap) bubbleUp(i int) {
	j := parent(i)
	if j == -1 {
		return
	}
	if h.heap[i].lessThan(h.heap[j]) {
		h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
		h.bubbleUp(j)
	}
}

func (h *fheap) placeMin() bool {
	if h.n == 0 {
		return false
	}

	if h.heap[0].alive {
		h.out.WriteString(h.heap[0].head + "\n")
		h.heap[0].pop()
	} else {
		h.heap[0] = h.heap[h.n-1]
		h.n--
	}
	h.bubbleDown(0)
	return true
}

func (h *fheap) bubbleDown(i int) {
	c := child(i)
	min := i
	for leftRight := 0; leftRight <= 1; leftRight++ {
		if c+leftRight < h.n {
			if h.heap[c+leftRight].lessThan(h.heap[min]) {
				min = c + leftRight
			}
		}
	}
	if min != i {
		h.heap[i], h.heap[min] = h.heap[min], h.heap[i]
		h.bubbleDown(min)
	}
}

//MergeChunks merges a list of sorted files into one.
func MergeChunks(files []string, out string, keyVal KeyValFunction) error {
	fOut, err := os.Create(out)
	if err != nil {
		return err
	}
	defer fOut.Close()
	writer := bufio.NewWriter(fOut)
	defer writer.Flush()
	h := &fheap{0, make([]*chunk, 0, len(files)), writer}

	for _, f := range files {
		fIn, err := os.Open(f)
		if err != nil {
			return err
		}
		defer fIn.Close()

		reader := bufio.NewScanner(fIn)

		c := &chunk{"", "", reader, true, keyVal}
		c.pop()
		h.insert(c)
	}

	bar := pb.StartNew(-1)
	for h.placeMin() {
		bar.Add(1)
	}
	for _, scanner := range h.heap {
		if scanner.reader.Err() != nil {
			return scanner.reader.Err()
		}
	}
	bar.Finish()
	return nil
}

//keySorter allows us to sort by keys, not just the whole line.
type keySorter struct {
	toSort         []string
	keyValFunction KeyValFunction
}

func (k *keySorter) Len() int {
	return len(k.toSort)
}

func (k *keySorter) Swap(i, j int) {
	k.toSort[i], k.toSort[j] = k.toSort[j], k.toSort[i]
}

func (k *keySorter) Less(i, j int) bool {
	kI, _ := k.keyValFunction(k.toSort[i])
	kJ, _ := k.keyValFunction(k.toSort[j])
	return kI < kJ
}

//sortIntoChunks produces individually sorted chunks of a file.
func sortIntoChunks(scanner *bufio.Scanner, tempPath string, nBytes int, keyVal KeyValFunction) error {
	bar := pb.StartNew(-1)
	bar.Set(pb.Bytes, true)

	fIndex := 0
	toSort := make([]string, 0)
	continueScanning := true
	for continueScanning {
		bytesRead := 0
		fOut, err := os.Create(path.Join(tempPath, fmt.Sprintf("f%d.tmp", fIndex)))
		if err != nil {
			return err
		}
		writer := bufio.NewWriterSize(fOut, nBytes/2)

		for bytesRead < nBytes/2 {
			continueScanning = scanner.Scan()
			if !continueScanning {
				if scanner.Err() != nil {
					return scanner.Err()
				}
				break
			}
			scannerBytes := len(scanner.Bytes())
			bar.Add(scannerBytes)
			bytesRead += scannerBytes
			toSort = append(toSort, scanner.Text())
		}

		sort.Sort(&keySorter{toSort, keyVal})
		for _, line := range toSort {
			writer.WriteString(line + "\n")
		}
		writer.Flush()
		fIndex++
		toSort = toSort[:0]
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	bar.Finish()
	return nil
}

//ExternalSort sorts the lines of in and writes them to out. Attempts to use no
//more than nBytes of space.
func ExternalSort(in, out string, nBytes int, keyVal KeyValFunction) error {
	fIn, err := os.Open(in)
	if err != nil {
		return err
	}
	defer fIn.Close()

	scanner := bufio.NewScanner(fIn)
	buffer := make([]byte, nBytes/2, nBytes/2)
	scanner.Buffer(buffer, nBytes/2)

	path := filepath.Dir(out)
	tempPath, _ := ioutil.TempDir(path, "sort")
	defer os.RemoveAll(tempPath)

	err = sortIntoChunks(scanner, tempPath, nBytes, keyVal)
	if err != nil {
		return err
	}
	chunkPaths, _ := filepath.Glob(tempPath + "/*")
	return MergeChunks(chunkPaths, out, keyVal)
}
