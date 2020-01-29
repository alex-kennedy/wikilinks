package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//ResolveRedirects resolves the redirect page titles.
type ResolveRedirects struct{}

//Run downloads the file.
func (t *ResolveRedirects) Run() error {
	log.Println("Resolving redirects...")
	pageRedirect := viper.GetString("page_redirect_sorted")
	resolved := viper.GetString("page_redirect_resolved")
	redirect := viper.GetString("redirect_sorted")
	redirectIndex := viper.GetString("redirect_index")
	pageDirect := viper.GetString("page_direct_sorted")
	pageDirectIndex := viper.GetString("page_direct_Index")
	bytesPerBuffer := viper.GetInt("bytes") / 10 //Arbitrary
	return lib.ResolveRedirects(pageRedirect, resolved, redirect, redirectIndex,
		pageDirect, pageDirectIndex, bytesPerBuffer)
}

//Done checks if the resolution completed successfully.
func (t *ResolveRedirects) Done() bool {
	return lib.CheckExists(viper.GetString("page_redirect_resolved"))
}

//Cleanup removes partial files on a failed resolution run.
func (t *ResolveRedirects) Cleanup() error {
	return lib.CleanupFile(viper.GetString("page_redirect_resolved"))
}

//Deps returns the dependencies of this task.
func (t *ResolveRedirects) Deps() []Task {
	return []Task{&IndexPageDirect{}, &IndexRedirect{}, &SortPageRedirect{}}
}
