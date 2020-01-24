package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//UnzipPagelinks extracts the the .sql.gz file to .sql
type UnzipPagelinks struct{}

//Run extracts the file
func (t *UnzipPagelinks) Run() error {
	log.Println("Unzipping pagelinks table...")
	inPath := viper.GetString("pagelinks_sql_gz")
	outPath := viper.GetString("pagelinks_sql")
	return lib.UnzipGzFile(inPath, outPath)
}

//Done checks if the unzip completed successfully.
func (t *UnzipPagelinks) Done() bool {
	return lib.CheckExists(viper.GetString("pagelinks_sql"))
}

//Cleanup removes partial files on a failed unzip.
func (t *UnzipPagelinks) Cleanup() error {
	return lib.CleanupFile(viper.GetString("pagelinks_sql"))
}

//Deps returns the dependencies of this task.
func (t *UnzipPagelinks) Deps() []Task {
	return []Task{&DownloadPagelinks{}}
}
