package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//SortPagelinks sorts the resolved pagelinks.
type SortPagelinks struct{}

//Run sorts the file
func (t *SortPagelinks) Run() error {
	log.Println("Sorting pagelinks...")
	inPath := viper.GetString("pagelinks_resolved")
	outPath := viper.GetString("pagelinks_sorted")
	nBytes := viper.GetInt("bytes")
	return lib.ExternalSort(inPath, outPath, nBytes, lib.KeyValFirstComma)
}

//Done checks if the sort completed successfully.
func (t *SortPagelinks) Done() bool {
	return lib.CheckExists(viper.GetString("pagelinks_sorted"))
}

//Cleanup removes partial files on a failed sort.
func (t *SortPagelinks) Cleanup() error {
	return lib.CleanupFile(viper.GetString("pagelinks_sorted"))
}

//Deps returns the dependencies of this task.
func (t *SortPagelinks) Deps() []Task {
	return []Task{&ResolvePagelinks{}}
}
