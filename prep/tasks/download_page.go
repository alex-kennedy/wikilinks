package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//DownloadPage downloads page.sql.gz.
type DownloadPage struct{}

//Run downloads the file.
func (t *DownloadPage) Run() error {
	log.Println("Downloading page table...")
	fileName := "enwiki-" + viper.GetString("date") + "-page.sql.gz"
	siteURL := viper.GetString("site_url") + viper.GetString("date") + "/"
	outPath := viper.GetString("page_sql_gz")
	return lib.DownloadWikiFile(fileName, siteURL, outPath)
}

//Done checks if the download completed successfully.
func (t *DownloadPage) Done() bool {
	return lib.CheckExists(viper.GetString("page_sql_gz"))
}

//Cleanup removes partial files on a failed download.
func (t *DownloadPage) Cleanup() error {
	return lib.CleanupFile(viper.GetString("page_sql_gz"))
}

//Deps returns the dependencies of this task.
func (t *DownloadPage) Deps() []Task {
	return []Task{&CreateFolders{}}
}
