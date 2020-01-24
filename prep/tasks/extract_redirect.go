package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//ExtractRedirect extracts data from redirect.sql.
//https://www.mediawiki.org/wiki/Manual:Redirect_table
type ExtractRedirect struct{}

//Run does the extraction.
func (t *ExtractRedirect) Run() error {
	log.Println("Extracting redirect...")
	inPath := viper.GetString("redirect_sql")
	outPath := viper.GetString("redirect")
	indices := []int{2, 0}
	fieldsPerRecord := 5
	return lib.ExtractTable(inPath, outPath, indices, fieldsPerRecord, keepRedirect)
}

//Done checks if the extraction completed successfully.
func (t *ExtractRedirect) Done() bool {
	return lib.CheckExists(viper.GetString("redirect"))
}

//Cleanup removes partial files on a failed extraction.
func (t *ExtractRedirect) Cleanup() error {
	return lib.CleanupFile(viper.GetString("redirect"))
}

//Deps returns the dependencies of this task.
func (t *ExtractRedirect) Deps() []Task {
	return []Task{&UnzipRedirect{}}
}

//keepRedirect if rd_namespace == "0"
func keepRedirect(record []string) bool {
	return record[1] == "0"
}
