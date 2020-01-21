package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//DownloadPagelinks downloads pagelinks.sql.gz.
type DownloadPagelinks struct{}

//Run downloads the file.
func (t *DownloadPagelinks) Run() error {
	log.Println("Downloading pagelinks table...")
	fileName := "enwiki-" + viper.GetString("date") + "-pagelinks.sql.gz"
	siteURL := viper.GetString("site_url") + viper.GetString("date") + "/"
	outPath := viper.GetString("pagelinks_sql_gz")
	return lib.DownloadWikiFile(fileName, siteURL, outPath)
}

//Done checks if the download completed successfully.
func (t *DownloadPagelinks) Done() bool {
	return lib.CheckExists(viper.GetString("pagelinks_sql_gz"))
}

//Cleanup removes partial files on a failed download.
func (t *DownloadPagelinks) Cleanup() error {
	return lib.CleanupFile(viper.GetString("pagelinks_sql_gz"))
}

//Deps returns the dependencies of this task.
func (t *DownloadPagelinks) Deps() []Task {
	return []Task{&CreateFolders{}}
}
