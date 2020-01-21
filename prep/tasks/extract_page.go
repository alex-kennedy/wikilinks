package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//ExtractPage extracts data from page.sql.
//https://www.mediawiki.org/wiki/Manual:Page_table
type ExtractPage struct{}

//Run does the extraction.
func (t *ExtractPage) Run() error {
	log.Println("Extracting page...")
	inPath := viper.GetString("page_sql")
	outPath := viper.GetString("page")
	indices := []int{2, 0, 1, 4}
	fieldsPerRecord := 13
	return lib.ExtractTable(inPath, outPath, indices, fieldsPerRecord)
}

//Done checks if the extraction completed successfully.
func (t *ExtractPage) Done() bool {
	return lib.CheckExists(viper.GetString("page"))
}

//Cleanup removes partial files on a failed extraction.
func (t *ExtractPage) Cleanup() error {
	return lib.CleanupFile(viper.GetString("page"))
}

//Deps returns the dependencies of this task.
func (t *ExtractPage) Deps() []Task {
	return []Task{&DownloadPage{}}
}
