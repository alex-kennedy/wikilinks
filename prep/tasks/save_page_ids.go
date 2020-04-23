package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//SavePageIDs extracts all the page IDs from page_direct (which is to say, valid
//page IDs), sorts them, and writes them to a file.
type SavePageIDs struct{}

//Run gets the page IDs.
func (t *SavePageIDs) Run() error {
	log.Println("Getting page IDs...")
	inPath := viper.GetString("page_direct")
	outPath := viper.GetString("page_ids")
	return lib.SavePageIDs(inPath, outPath)
}

//Done checks if the extraction completed successfully.
func (t *SavePageIDs) Done() bool {
	return lib.CheckExists(viper.GetString("page_ids"))
}

//Cleanup removes partial files on a failed extraction.
func (t *SavePageIDs) Cleanup() error {
	return lib.CleanupFile(viper.GetString("page_ids"))
}

//Deps returns the dependencies of this task.
func (t *SavePageIDs) Deps() []Task {
	return []Task{&SortPageDirect{}}
}
