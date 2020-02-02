package lib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/cheggaaa/pb"
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

	pb := pb.StartNew(-1)
	defer pb.Finish()

	for pageRedirectScanner.Scan() {
		line := pageRedirectScanner.Text()
		kA, vA := KeyValLastComma(line)
		if vA == "" {
			continue
		}

		// Phase A: pageRedirect[redirect page ID] -> redirect[title of dest page]
		_, kB, err := redirectSearcher.Search(vA)
		if err != nil {
			status.failedRedirectLooksups++
			continue
		}
		status.successfulRedirectLookups++

		// Phase B: redirect[title of dest page] -> pageDirect[true page ID]
		_, vB, err := pageDirectSearcher.Search(kB)
		if err != nil {
			status.failedPageLookups++
			continue
		}
		status.successfulPageLookups++

		resolvedWriter.WriteString(kA + "," + vB + "\n")
		pb.Add(1)
	}

	status.printStatus()
	return nil
}

//ResolvePagelinks turns pagelinks titles into IDs and saves them as base36 IDs.
func ResolvePagelinks(pageMerged, pageMergedIndex, pagelinks, out string,
	bytesPerBuffer int) error {
	successful, failed := 0, 0

	pagelinksFile, err := os.Open(pagelinks)
	if err != nil {
		return err
	}
	defer pagelinksFile.Close()
	pagelinksScanner := bufio.NewScanner(pagelinksFile)
	scannerBuffer := make([]byte, bytesPerBuffer)
	pagelinksScanner.Buffer(scannerBuffer, bytesPerBuffer)

	outFile, err := os.Create(out)
	if err != nil {
		return err
	}
	defer outFile.Close()
	outWriter := bufio.NewWriterSize(outFile, bytesPerBuffer)
	defer outWriter.Flush()

	pageSearcher, err := NewBinarySearcher(pageMerged, pageMergedIndex, KeyValFirstComma)
	if err != nil {
		return err
	}

	pb := pb.StartNew(-1)
	defer pb.Finish()

	for pagelinksScanner.Scan() {
		line := pagelinksScanner.Text()
		key, title := KeyValFirstComma(line)
		if title == "" {
			failed++
			continue
		}

		_, titleID, err := pageSearcher.Search(title)
		if err != nil {
			failed++
			fmt.Println(title)
			continue
		}

		keyInt, err := strconv.Atoi(key)
		if err != nil {
			failed++
			log.Printf("'%s' couldn't be parsed to int", key)
			continue
		}
		titleIDInt, err := strconv.Atoi(titleID)
		if err != nil {
			failed++
			log.Printf("'%s' couldn't be parsed to int", titleID)
			continue
		}

		successful++
		outWriter.WriteString(strconv.FormatInt(int64(keyInt), 36) + "," +
			strconv.FormatInt(int64(titleIDInt), 36) + "\n")
		pb.Add(1)
	}

	log.Printf("%d failed, %d lines succeeded", failed, successful)
	return nil
}
