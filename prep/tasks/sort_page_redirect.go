package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//SortPageRedirect sorts the redirected parts of the page file.
type SortPageRedirect struct{}

//Run sorts the file
func (t *SortPageRedirect) Run() error {
	log.Println("Sorting page redirect...")
	inPath := viper.GetString("page_redirect")
	outPath := viper.GetString("page_redirect_sorted")
	nBytes := viper.GetInt("bytes")
	return lib.ExternalSort(inPath, outPath, nBytes, lib.KeyValLastComma)
}

//Done checks if the sort completed successfully.
func (t *SortPageRedirect) Done() bool {
	return lib.CheckExists(viper.GetString("page_redirect_sorted"))
}

//Cleanup removes partial files on a failed sort.
func (t *SortPageRedirect) Cleanup() error {
	return lib.CleanupFile(viper.GetString("page_redirect_sorted"))
}

//Deps returns the dependencies of this task.
func (t *SortPageRedirect) Deps() []Task {
	return []Task{&ExtractPageRedirect{}}
}
