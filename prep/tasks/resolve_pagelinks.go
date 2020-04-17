package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//ResolvePagelinks resolves the pagelinks page titles.
type ResolvePagelinks struct{}

//Run downloads the file.
func (t *ResolvePagelinks) Run() error {
	log.Println("Resolving pagelinks...")
	pageMerged := viper.GetString("page_merged")
	pageMergedIndex := viper.GetString("page_merged_index")
	pagelinks := viper.GetString("pagelinks")
	out := viper.GetString("pagelinks_resolved")
	return lib.ResolvePagelinks(pageMerged, pageMergedIndex, pagelinks, out)
}

//Done checks if the resolution completed successfully.
func (t *ResolvePagelinks) Done() bool {
	return lib.CheckExists(viper.GetString("pagelinks_resolved"))
}

//Cleanup removes partial files on a failed resolution run.
func (t *ResolvePagelinks) Cleanup() error {
	return lib.CleanupFile(viper.GetString("pagelinks_resolved"))
}

//Deps returns the dependencies of this task.
func (t *ResolvePagelinks) Deps() []Task {
	return []Task{&ExtractPagelinks{}, &IndexPageMerged{}}
}
