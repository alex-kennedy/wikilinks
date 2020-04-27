package lib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/cheggaaa/pb/v3"
)

//CreatePagelinksLGL takes the pagelinks_pivoted file and outputs it in the
//format expected by the LGL implementation. Performs de-duplication as the
//implementation can't handle bi-directional links. For the purposes of
//colouring them later on, the lower page ID is considered canonical. I'll
//encode them in base-36 for the big ol Wikipedia one.
func CreatePagelinksLGL(pagelinksPivotedName, pageIDsName, outName string) error {
	outFile, err := os.Create(outName)
	if err != nil {
		return err
	}
	out := bufio.NewWriter(outFile)
	defer out.Flush()

	pagelinks, err := NewPagelinksPivotedInMemory(pagelinksPivotedName, pageIDsName)
	if err != nil {
		return err
	}

	bar := pb.StartNew(pagelinks.GetNumberOfPages())
	defer bar.Finish()
	numberOfBis := 0
	firstWrite := true
	for from := 0; from < pagelinks.GetNumberOfPages(); from++ {
		bar.Add(1)

		toSlice := pagelinks.GetLinks(from)
		if len(toSlice) == 0 {
			continue
		}

		//Header lines.
		//Prevents any trailing new lines.
		if firstWrite {
			fmt.Fprintf(out, "# %s", strconv.FormatUint(uint64(from), 36))
			firstWrite = false
		} else {
			fmt.Fprintf(out, "\n# %s", strconv.FormatUint(uint64(from), 36))
		}

		for _, to := range toSlice {
			//If from < to, then this is the canonical link so it is always written.
			//If it is not canonical, only write if the reverse doesn't exist.
			if uint32(from) < to || !pagelinks.LinkExists(to, uint32(from)) {
				fmt.Fprintf(out, "\n%s", strconv.FormatUint(uint64(to), 36))
			} else {
				numberOfBis++
			}
		}
	}
	log.Printf("Done. There were %d bi-directional links", numberOfBis)
	return nil
}
