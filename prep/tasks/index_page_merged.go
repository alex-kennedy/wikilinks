package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//IndexPageMerged indexes the merged, resolved page table.
type IndexPageMerged struct{}

//Run indexes the file
func (t *IndexPageMerged) Run() error {
	log.Println("Indexing page merged...")
	inPath := viper.GetString("page_merged")
	outPath := viper.GetString("page_merged_index")
	return lib.IndexFile(inPath, outPath)
}

//Done checks if the indexing completed successfully.
func (t *IndexPageMerged) Done() bool {
	return lib.CheckExists(viper.GetString("page_merged_index"))
}

//Cleanup removes partial files on a failed indexing.
func (t *IndexPageMerged) Cleanup() error {
	return lib.CleanupFile(viper.GetString("page_merged_index"))
}

//Deps returns the dependencies of this task.
func (t *IndexPageMerged) Deps() []Task {
	return []Task{&MergePage{}}
}
