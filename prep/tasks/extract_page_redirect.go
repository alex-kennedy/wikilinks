package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//ExtractPageRedirect extracts data from page.sql that are redirects.
//https://www.mediawiki.org/wiki/Manual:Page_table
type ExtractPageRedirect struct{}

//Run does the extraction.
func (t *ExtractPageRedirect) Run() error {
	log.Println("Extracting page redirect...")
	inPath := viper.GetString("page_sql")
	outPath := viper.GetString("page_redirect")
	// Actual indices differ from the wiki - check on update
	indices := []int{2, 0}
	fieldsPerRecord := 13
	return lib.ExtractTable(inPath, outPath, indices, fieldsPerRecord, keepPageRedirect)
}

//Done checks if the extraction completed successfully.
func (t *ExtractPageRedirect) Done() bool {
	return lib.CheckExists(viper.GetString("page_redirect"))
}

//Cleanup removes partial files on a failed extraction.
func (t *ExtractPageRedirect) Cleanup() error {
	return lib.CleanupFile(viper.GetString("page_redirect"))
}

//Deps returns the dependencies of this task.
func (t *ExtractPageRedirect) Deps() []Task {
	return []Task{&UnzipPage{}}
}

//keepPage if page_namespace == "0" and page_is_redirect == "1"
func keepPageRedirect(record []string) bool {
	return record[1] == "0" && record[4] == "1"
}
