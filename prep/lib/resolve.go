package lib

import (
	"bufio"
	"github.com/cheggaaa/pb"
	"log"
	"os"
)

//ResolveRedirects resolves the page_redirect -> redirect -> page_direct
//problem.
func ResolveRedirects(pageRedirect, resolved, redirect, redirectIndex,
	pageDirect, pageDirectIndex string, bytesPerBuffer int) error {

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
			log.Printf("%s does not exist in redirect", vA)
			continue
		}

		// Phase B: redirect[title of dest page] -> pageDirect[true page ID]
		_, vB, err := pageDirectSearcher.Search(kB)
		if err != nil {
			log.Printf("%s does not exist in page_direct", kB)
			continue
		}

		resolvedWriter.WriteString(kA + "," + vB + "\n")
		pb.Add(1)
	}

	return nil
}
