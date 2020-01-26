package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//IndexPageDirect indexes the direct parts of the page file.
type IndexPageDirect struct{}

//Run indexes the file
func (t *IndexPageDirect) Run() error {
	log.Println("Indexing page direct...")
	inPath := viper.GetString("page_direct_sorted")
	outPath := viper.GetString("page_direct_index")
	return lib.IndexFile(inPath, outPath)
}

//Done checks if the unzip completed successfully.
func (t *IndexPageDirect) Done() bool {
	return lib.CheckExists(viper.GetString("page_direct_index"))
}

//Cleanup removes partial files on a failed unzip.
func (t *IndexPageDirect) Cleanup() error {
	return lib.CleanupFile(viper.GetString("page_direct_index"))
}

//Deps returns the dependencies of this task.
func (t *IndexPageDirect) Deps() []Task {
	return []Task{&SortPageDirect{}}
}
