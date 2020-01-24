package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//UnzipRedirect extracts the the .sql.gz file to .sql
type UnzipRedirect struct{}

//Run extracts the file
func (t *UnzipRedirect) Run() error {
	log.Println("Unzipping redirect table...")
	inPath := viper.GetString("redirect_sql_gz")
	outPath := viper.GetString("redirect_sql")
	return lib.UnzipGzFile(inPath, outPath)
}

//Done checks if the unzip completed successfully.
func (t *UnzipRedirect) Done() bool {
	return lib.CheckExists(viper.GetString("redirect_sql"))
}

//Cleanup removes partial files on a failed unzip.
func (t *UnzipRedirect) Cleanup() error {
	return lib.CleanupFile(viper.GetString("redirect_sql"))
}

//Deps returns the dependencies of this task.
func (t *UnzipRedirect) Deps() []Task {
	return []Task{&DownloadRedirect{}}
}
