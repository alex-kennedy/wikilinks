package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//IndexRedirect indexes the sorted redirect table.
type IndexRedirect struct{}

//Run indexes the file
func (t *IndexRedirect) Run() error {
	log.Println("Indexing redirect...")
	inPath := viper.GetString("redirect_sorted")
	outPath := viper.GetString("redirect_index")
	return lib.IndexFile(inPath, outPath)
}

//Done checks if the unzip completed successfully.
func (t *IndexRedirect) Done() bool {
	return lib.CheckExists(viper.GetString("redirect_index"))
}

//Cleanup removes partial files on a failed unzip.
func (t *IndexRedirect) Cleanup() error {
	return lib.CleanupFile(viper.GetString("redirect_index"))
}

//Deps returns the dependencies of this task.
func (t *IndexRedirect) Deps() []Task {
	return []Task{&SortRedirect{}}
}
