package tasks

import (
	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//DoPageRank ensures the data directory exists.
type DoPageRank struct{}

//Run creates directory if it doesn't exist
func (t *DoPageRank) Run() error {
	backlinksName := viper.GetString("backlinks")
	pagelinksName := viper.GetString("pagelinks_pivoted")
	pageIDsName := viper.GetString("page_ids")
	pageranksName := viper.GetString("pageranks")

	pagelinksCounts, err := lib.CountPagelinks(pagelinksName, pageIDsName)
	if err != nil {
		return err
	}

	backlinks, err := lib.NewPagelinksPivotedInMemory(backlinksName, pageIDsName)
	if err != nil {
		return err
	}
	ranker := lib.NewPageRanker(backlinks, pagelinksCounts, 60, 0.85)
	ranker.Rank()
	ranker.OutputToFile(pageranksName)
	return nil
}

//Done checks that the pageranks file exists.
func (t *DoPageRank) Done() bool {
	return lib.CheckExists(viper.GetString("pageranks"))
}

//Cleanup removes any leftover files.
func (t *DoPageRank) Cleanup() error {
	return lib.CleanupFile(viper.GetString("pageranks"))
}

//Deps returns the dependencies of this task.
func (t *DoPageRank) Deps() []Task {
	return []Task{&SaveBacklinks{}}
}
