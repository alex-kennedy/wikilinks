package lib

import (
	"bufio"
	"fmt"
	"github.com/cheggaaa/pb"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
)

//chunk represents one sorted file
type chunk struct {
	head   string
	reader *bufio.Scanner
	alive  bool
}

func (c *chunk) pop() {
	c.alive = c.reader.Scan()
	c.head = c.reader.Text()
}

func (c *chunk) lessThan(d *chunk) bool {
	return c.head < d.head
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

func mergeChunks(files []string, out string, bufferSize int) error {
	fOut, err := os.Create(out)
	if err != nil {
		return err
	}
	defer fOut.Close()
	writer := bufio.NewWriterSize(fOut, bufferSize)
	defer writer.Flush()
	h := &fheap{0, make([]*chunk, 0, len(files)), writer}

	for _, f := range files {
		fIn, err := os.Open(f)
		if err != nil {
			return err
		}
		defer fIn.Close()

		reader := bufio.NewScanner(fIn)
		buffer := make([]byte, bufferSize)
		reader.Buffer(buffer, bufferSize)

		c := &chunk{"", reader, true}
		c.pop()
		h.insert(c)
	}

	log.Println("Merging...")
	bar := pb.StartNew(-1)
	for h.placeMin() {
		bar.Add(1)
	}
	return nil
}

//sortIntoChunks produces individually sorted chunks of a file.
func sortIntoChunks(scanner *bufio.Scanner, tempPath string, nBytes int) error {
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

		sort.Strings(toSort)
		for _, line := range toSort {
			writer.WriteString(line + "\n")
		}
		writer.Flush()
		fIndex++
		toSort = toSort[:0]
	}

	return nil
}

//ExternalSort sorts the lines of in and writes them to out. Attempts to use no
//more than nBytes of space.
func ExternalSort(in, out string, nBytes int) error {
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

	sortIntoChunks(scanner, tempPath, nBytes)
	chunkPaths, _ := filepath.Glob(tempPath + "/*")
	return mergeChunks(chunkPaths, out, nBytes/(len(chunkPaths)+1))
}
