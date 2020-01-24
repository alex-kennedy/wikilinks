package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//ExtractPageDirect extracts data from page.sql that are not redirects.
//https://www.mediawiki.org/wiki/Manual:Page_table
type ExtractPageDirect struct{}

//Run does the extraction.
func (t *ExtractPageDirect) Run() error {
	log.Println("Extracting page direct...")
	inPath := viper.GetString("page_sql")
	outPath := viper.GetString("page_direct")
	indices := []int{2, 0}
	fieldsPerRecord := 13
	return lib.ExtractTable(inPath, outPath, indices, fieldsPerRecord, keepPageDirect)
}

//Done checks if the extraction completed successfully.
func (t *ExtractPageDirect) Done() bool {
	return lib.CheckExists(viper.GetString("page_direct"))
}

//Cleanup removes partial files on a failed extraction.
func (t *ExtractPageDirect) Cleanup() error {
	return lib.CleanupFile(viper.GetString("page_direct"))
}

//Deps returns the dependencies of this task.
func (t *ExtractPageDirect) Deps() []Task {
	return []Task{&UnzipPage{}}
}

//keepPage if page_namespace == "0" and page_is_redirect == "0"
func keepPageDirect(record []string) bool {
	return record[1] == "0" && record[4] == "0"
}
