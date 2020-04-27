package lib

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

//PageRanker determines the PageRank of each page.
//https://en.wikipedia.org/wiki/PageRank
type PageRanker struct {
	Backlinks       PagelinksPivoted
	PagelinksCounts []float64
	MaxIterations   int
	DampingFactor   float64
	numberOfPages   int
	ranks           []float64
	ranksLast       []float64
	norms           []float64
	logAndScaled    []float64
}

//Rank runs the PageRank algorithm.
func (p *PageRanker) Rank() {
	//Initialise
	initValue := 1 / float64(p.numberOfPages)
	for i := range p.ranks {
		p.ranks[i] = initValue
	}
	p.norms = append(p.norms, 1.0)
	log.Printf("Iteration 0: %f", p.norms[0])
	for it := 1; it <= p.MaxIterations; it++ {
		p.doIteration()
		log.Printf("Iterations %d: %e (diff of %e)", it, p.norms[it], math.Abs(p.norms[it]-p.norms[it-1]))
		if p.norms[it] < 1e-10 {
			log.Printf("Tolerance reached. Exiting after %d iterations", it)
			break
		}
	}
	p.LogAndScale()
}

//LogAndScale performs a natural logarithm on the last calculated set of ranks,
//and scales them to the [0, 1].
func (p *PageRanker) LogAndScale() {
	minRank, maxRank := 0.0, 0.0
	for i := 0; i < p.numberOfPages; i++ {
		p.logAndScaled[i] = math.Log(p.ranks[i])
		if p.logAndScaled[i] < minRank {
			minRank = p.logAndScaled[i]
		}
		if p.logAndScaled[i] > maxRank {
			maxRank = p.logAndScaled[i]
		}
	}
	minMaxDiff := maxRank - minRank
	for i := 0; i < p.numberOfPages; i++ {
		p.logAndScaled[i] = (p.logAndScaled[i] - minRank) / minMaxDiff
	}
}

//OutputToFile writes the LogAndScale()d ranks to a file.
func (p *PageRanker) OutputToFile(name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	out := bufio.NewWriter(file)
	defer file.Close()
	defer out.Flush()

	for i := 0; i < p.numberOfPages; i++ {
		fmt.Fprintln(out, strconv.FormatFloat(p.logAndScaled[i], 'f', -1, 64))
	}
	return nil
}

//doIteration runs one iteration of the PageRank algorithm. Also calculates
//something like the norm of the vector (sum of absolute values) to identify
//convergence.
func (p *PageRanker) doIteration() {
	p.ranks, p.ranksLast = p.ranksLast, p.ranks
	constant := (1.0 - p.DampingFactor) / float64(p.numberOfPages)
	thisNorm := 0.0
	for i := range p.ranks {
		rank := 0.0
		for _, j := range p.Backlinks.GetLinks(i) {
			if p.PagelinksCounts[j] > 1e-10 {
				rank += p.ranksLast[j] / p.PagelinksCounts[j]
			}
		}
		rank *= p.DampingFactor
		rank += constant
		p.ranks[i] = rank
		thisNorm += math.Abs(rank - p.ranksLast[i])
	}
	p.norms = append(p.norms, thisNorm/float64(p.numberOfPages))
}

//NewPageRanker creates a PageRanker object to run the PageRank algorithm.
func NewPageRanker(backlinks PagelinksPivoted, pagelinksCounts []float64,
	maxIterations int, dampingFactor float64) *PageRanker {
	numberOfPages := backlinks.GetNumberOfPages()
	return &PageRanker{
		Backlinks:       backlinks,
		PagelinksCounts: pagelinksCounts,
		MaxIterations:   maxIterations,
		DampingFactor:   dampingFactor,
		numberOfPages:   numberOfPages,
		ranks:           make([]float64, numberOfPages),
		ranksLast:       make([]float64, numberOfPages),
		logAndScaled:    make([]float64, numberOfPages),
	}
}

//CountPagelinks counts the number of outbound links in each page. Values used
//as part of the PageRank algorithm.
func CountPagelinks(pagelinksPivotedName, pageIDsName string) ([]float64, error) {
	log.Println("Counting pagelinks...")
	pageIDs, err := LoadPageIDs(pageIDsName)
	if err != nil {
		return nil, err
	}
	numberOfPages := len(pageIDs)
	pageIDs = nil

	pagelinksPivotedFile, err := os.Open(pagelinksPivotedName)
	if err != nil {
		return nil, err
	}
	pagelinks := bufio.NewScanner(pagelinksPivotedFile)
	defer pagelinksPivotedFile.Close()
	bar := NewProgressBarFileSize(pagelinksPivotedFile)
	defer bar.Finish()

	pagelinkCounts := make([]float64, numberOfPages)
	var line string
	var lineIDs []string
	for pagelinks.Scan() {
		line = pagelinks.Text()
		bar.Add(len(line) + 1)
		lineIDs = strings.Split(line, ",")
		if len(lineIDs) > 1 {
			pagelinkCounts[MustParseBase36(lineIDs[0])] = float64(len(lineIDs) - 1)
		}
	}
	if pagelinks.Err() != nil {
		return nil, pagelinks.Err()
	}
	return pagelinkCounts, nil
}
