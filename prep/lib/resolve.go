package lib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

//Status stores stats about the resolution process.
type Status struct {
	successfulRedirectLookups int
	failedRedirectLooksups    int
	successfulPageLookups     int
	failedPageLookups         int
}

func (s *Status) printStatus() {
	log.Printf("redirect %d failed, %d succeeded", s.failedRedirectLooksups, s.successfulRedirectLookups)
	log.Printf("page: %d failed, %d succeeded", s.failedPageLookups, s.successfulPageLookups)
}

//ResolveRedirects resolves the page_redirect -> redirect -> page_direct
//problem.
func ResolveRedirects(pageRedirect, resolved, redirect, redirectIndex,
	pageDirect, pageDirectIndex string, bytesPerBuffer int) error {
	status := Status{}
	status.printStatus()

	pageRedirectFile, err := os.Open(pageRedirect)
	if err != nil {
		return err
	}
	defer pageRedirectFile.Close()
	pageRedirectScanner := bufio.NewScanner(pageRedirectFile)
	scannerBuffer := make([]byte, bytesPerBuffer)
	pageRedirectScanner.Buffer(scannerBuffer, bytesPerBuffer)

	resolvedFile, err := os.Create(resolved)
	if err != nil {
		return err
	}
	defer resolvedFile.Close()
	resolvedWriter := bufio.NewWriterSize(resolvedFile, bytesPerBuffer)
	defer resolvedWriter.Flush()

	redirectSearcher, err := NewBinarySearcher(redirect, redirectIndex, KeyValFirstComma)
	if err != nil {
		return err
	}
	pageDirectSearcher, err := NewBinarySearcher(pageDirect, pageDirectIndex, KeyValLastComma)
	if err != nil {
		return err
	}

	pb := NewProgressBarFileSize(pageRedirectFile)
	defer pb.Finish()

	for pageRedirectScanner.Scan() {
		line := pageRedirectScanner.Text()
		kA, vA := KeyValLastComma(line)
		if vA == "" {
			continue
		}

		// Phase A: pageRedirect[redirect page ID] -> redirect[title of dest page]
		kB, err := redirectSearcher.Search(vA)
		if err != nil {
			status.failedRedirectLooksups++
			continue
		}
		status.successfulRedirectLookups++

		// Phase B: redirect[title of dest page] -> pageDirect[true page ID]
		vB, err := pageDirectSearcher.Search(kB)
		if err != nil {
			status.failedPageLookups++
			continue
		}
		status.successfulPageLookups++

		resolvedWriter.WriteString(kA + "," + vB + "\n")
		pb.Add(len(line) + 1)
	}

	status.printStatus()
	return nil
}

//ResolvePagelinks turns pagelinks titles into IDs and saves them as base36 IDs
//(to reduce disk space). A page link is only resolved if the page_id from which
//the link originates is in the page_merged file. This is because redirected
//pages are stored as a row in the pagelinks table. We have already resolved
//this problem by bypassing redirect links, so they are skipped.
func ResolvePagelinks(pageMerged, pagelinks, out string) error {
	pagelinksFile, err := os.Open(pagelinks)
	if err != nil {
		return err
	}
	defer pagelinksFile.Close()
	pagelinksScanner := bufio.NewScanner(pagelinksFile)

	outFile, err := os.Create(out)
	if err != nil {
		return err
	}
	defer outFile.Close()
	outWriter := bufio.NewWriter(outFile)
	defer outWriter.Flush()

	log.Println("Loading page merged...")
	pageSearcher, err := NewStringToIntArraySearcher(pageMerged)
	if err != nil {
		return err
	}
	//Gets sorted list of values in the pageSearcher (pageMerged)
	pageSearcherValues := make([]int, len(pageSearcher.values))
	copy(pageSearcherValues, pageSearcher.values)
	sort.Ints(pageSearcherValues)
	log.Println("Done.")

	bar := NewProgressBarFileSize(pagelinksFile)
	defer bar.Finish()

	successful, failed, redirects := 0, 0, 0
	var titleID, keyInt, i int
	var key, title, line string
	for pagelinksScanner.Scan() {
		line = pagelinksScanner.Text()
		key, title = KeyValFirstComma(line)
		if title == "" {
			failed++
			continue
		}
		keyInt, err = strconv.Atoi(key)
		if err != nil {
			failed++
			log.Printf("'%s' couldn't be parsed to int", key)
			continue
		}

		i = sort.SearchInts(pageSearcherValues, keyInt)
		if i >= len(pageSearcherValues) || pageSearcherValues[i] != keyInt {
			//Refers to a redirect, skipping
			redirects++
			continue
		}

		titleID = pageSearcher.Search(title)
		if titleID == -1 {
			failed++
			continue
		}

		//Removes self-links
		if keyInt == titleID {
			failed++
			continue
		}

		successful++
		outWriter.WriteString(fmt.Sprintf("%s,%s\n",
			strconv.FormatInt(int64(keyInt), 36),
			strconv.FormatInt(int64(titleID), 36)))
		bar.Add(len(line) + 1)
	}

	log.Printf("%d failed, %d succeeded, %d redirected", failed, successful,
		redirects)
	return nil
}
