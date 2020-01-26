package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//SortPageDirect sorts the direct parts of the page file.
type SortPageDirect struct{}

//Run sorts the file
func (t *SortPageDirect) Run() error {
	log.Println("Sorting page direct...")
	inPath := viper.GetString("page_direct")
	outPath := viper.GetString("page_direct_sorted")
	nBytes := viper.GetInt("bytes")
	return lib.ExternalSort(inPath, outPath, nBytes)
}

//Done checks if the sort completed successfully.
func (t *SortPageDirect) Done() bool {
	return lib.CheckExists(viper.GetString("page_direct_sorted"))
}

//Cleanup removes partial files on a failed sort.
func (t *SortPageDirect) Cleanup() error {
	return lib.CleanupFile(viper.GetString("page_direct_sorted"))
}

//Deps returns the dependencies of this task.
func (t *SortPageDirect) Deps() []Task {
	return []Task{&ExtractPageDirect{}}
}
