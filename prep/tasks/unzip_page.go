package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//UnzipPage unzips the the .sql.gz file to .sql
type UnzipPage struct{}

//Run extracts the file
func (t *UnzipPage) Run() error {
	log.Println("Unzipping page table...")
	inPath := viper.GetString("page_sql_gz")
	outPath := viper.GetString("page_sql")
	return lib.UnzipGzFile(inPath, outPath)
}

//Done checks if the unzip completed successfully.
func (t *UnzipPage) Done() bool {
	return lib.CheckExists(viper.GetString("page_sql"))
}

//Cleanup removes partial files on a failed unzip.
func (t *UnzipPage) Cleanup() error {
	return lib.CleanupFile(viper.GetString("page_sql"))
}

//Deps returns the dependencies of this task.
func (t *UnzipPage) Deps() []Task {
	return []Task{&DownloadPage{}}
}
