package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//DownloadRedirect downloads redirect.sql.gz.
type DownloadRedirect struct{}

//Run downloads the file.
func (t *DownloadRedirect) Run() error {
	log.Println("Downloading redirect table...")
	fileName := "enwiki-" + viper.GetString("date") + "-redirect.sql.gz"
	siteURL := viper.GetString("site_url") + viper.GetString("date") + "/"
	outPath := viper.GetString("redirect_sql_gz")
	return lib.DownloadWikiFile(fileName, siteURL, outPath)
}

//Done checks if the download completed successfully.
func (t *DownloadRedirect) Done() bool {
	return lib.CheckExists(viper.GetString("redirect_sql_gz"))
}

//Cleanup removes partial files on a failed download.
func (t *DownloadRedirect) Cleanup() error {
	return lib.CleanupFile(viper.GetString("redirect_sql_gz"))
}

//Deps returns the dependencies of this task.
func (t *DownloadRedirect) Deps() []Task {
	return []Task{&CreateFolders{}}
}
