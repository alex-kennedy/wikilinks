package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//ExtractPagelinks extracts data from pagelinks.sql.
//https://www.mediawiki.org/wiki/Manual:Pagelinks_table
type ExtractPagelinks struct{}

//Run does the extraction.
func (t *ExtractPagelinks) Run() error {
	log.Println("Extracting pagelinks...")
	inPath := viper.GetString("pagelinks_sql")
	outPath := viper.GetString("pagelinks")
	indices := []int{0, 2}
	fieldsPerRecord := 4
	return lib.ExtractTable(inPath, outPath, indices, fieldsPerRecord)
}

//Done checks if the extraction completed successfully.
func (t *ExtractPagelinks) Done() bool {
	return lib.CheckExists(viper.GetString("pagelinks"))
}

//Cleanup removes partial files on a failed extraction.
func (t *ExtractPagelinks) Cleanup() error {
	return lib.CleanupFile(viper.GetString("pagelinks"))
}

//Deps returns the dependencies of this task.
func (t *ExtractPagelinks) Deps() []Task {
	return []Task{&DownloadPagelinks{}}
}
