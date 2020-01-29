package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//SortRedirect sorts the redirect table.
type SortRedirect struct{}

//Run sorts the file
func (t *SortRedirect) Run() error {
	log.Println("Sorting redirect...")
	inPath := viper.GetString("redirect")
	outPath := viper.GetString("redirect_sorted")
	nBytes := viper.GetInt("bytes")
	return lib.ExternalSort(inPath, outPath, nBytes, lib.KeyValFirstComma)
}

//Done checks if the sort completed successfully.
func (t *SortRedirect) Done() bool {
	return lib.CheckExists(viper.GetString("redirect_sorted"))
}

//Cleanup removes partial files on a failed sort.
func (t *SortRedirect) Cleanup() error {
	return lib.CleanupFile(viper.GetString("redirect_sorted"))
}

//Deps returns the dependencies of this task.
func (t *SortRedirect) Deps() []Task {
	return []Task{&ExtractRedirect{}}
}
