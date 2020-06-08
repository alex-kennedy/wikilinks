package tasks

import (
	"fmt"
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//DownloadRedirect downloads redirect.sql.gz.
type DownloadRedirect struct{}

//Run downloads the file.
func (t *DownloadRedirect) Run() error {
	log.Println("Downloading redirect table...")

	wikiName, err := lib.GetWikiNameFromURL(viper.GetString("site_url"))
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s-%s-redirect.sql.gz", wikiName, viper.GetString("date"))
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
