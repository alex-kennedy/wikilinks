package tasks

import (
	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//MergePage merges page_direct_sorted.csv and page_redirect_resolved.csv.
type MergePage struct{}

//Run merges the tables.
func (t *MergePage) Run() error {
	pageDirectSorted := viper.GetString("page_direct_sorted")
	pageRedirectSorted := viper.GetString("page_redirect_resolved")
	both := []string{pageDirectSorted, pageRedirectSorted}
	merged := viper.GetString("page_merged")
	bytes := viper.GetInt("bytes")
	return lib.MergeChunks(both, merged, bytes, lib.KeyValLastComma)
}

//Done checks if the download completed successfully.
func (t *MergePage) Done() bool {
	return lib.CheckExists(viper.GetString("page_merged"))
}

//Cleanup removes partial files on a failed download.
func (t *MergePage) Cleanup() error {
	return lib.CleanupFile(viper.GetString("page_merged"))
}

//Deps returns the dependencies of this task.
func (t *MergePage) Deps() []Task {
	return []Task{&ResolveRedirects{}}
}
