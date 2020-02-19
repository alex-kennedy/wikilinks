package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//IndexPagelinksPivoted indexes the pivoted pagelinks table.
type IndexPagelinksPivoted struct{}

//Run indexes the file
func (t *IndexPagelinksPivoted) Run() error {
	log.Println("Indexing pagelinks pivoted...")
	inPath := viper.GetString("pagelinks_pivoted")
	outPath := viper.GetString("pagelinks_pivoted_index")
	return lib.IndexFile(inPath, outPath)
}

//Done checks if the indexing completed successfully.
func (t *IndexPagelinksPivoted) Done() bool {
	return lib.CheckExists(viper.GetString("pagelinks_pivoted_index"))
}

//Cleanup removes partial files on a failed indexing.
func (t *IndexPagelinksPivoted) Cleanup() error {
	return lib.CleanupFile(viper.GetString("pagelinks_pivoted_index"))
}

//Deps returns the dependencies of this task.
func (t *IndexPagelinksPivoted) Deps() []Task {
	return []Task{&PivotPagelinks{}}
}
