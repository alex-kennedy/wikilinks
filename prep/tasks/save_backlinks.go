package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//SaveBacklinks inverts the pagelinks_pivoted table to get a backlinks table of
//the same form. For use in the PageRank implementation. Caution should be
//exercised, this loads all the pagelinks into memory.
type SaveBacklinks struct{}

//Run finds and saves the backlinks.
func (t *SaveBacklinks) Run() error {
	log.Println("Saving backlinks...")
	pagelinksPivoted := viper.GetString("pagelinks_pivoted")
	pageIDs := viper.GetString("page_ids")
	backlinks := viper.GetString("backlinks")
	return lib.SaveBacklinks(pagelinksPivoted, pageIDs, backlinks)
}

//Done checks if the task completed successfully.
func (t *SaveBacklinks) Done() bool {
	return lib.CheckExists(viper.GetString("backlinks"))
}

//Cleanup removes partial files on a failure.
func (t *SaveBacklinks) Cleanup() error {
	return lib.CleanupFile(viper.GetString("backlinks"))
}

//Deps returns the dependencies of this task.
func (t *SaveBacklinks) Deps() []Task {
	return []Task{&SortPageDirect{}}
}
