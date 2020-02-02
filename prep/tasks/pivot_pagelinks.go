package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//PivotPagelinks puts all the links in each file on the same line.
type PivotPagelinks struct{}

//Run downloads the file.
func (t *PivotPagelinks) Run() error {
	log.Println("Pivoting pagelinks...")
	pagelinksSorted := viper.GetString("pagelinks_sorted")
	pagelinksPivoted := viper.GetString("pagelinks_pivoted")
	bytesPerBuffer := viper.GetInt("bytes") / 10 //Arbitrary
	return lib.PivotFile(pagelinksSorted, pagelinksPivoted, bytesPerBuffer)
}

//Done checks if the pivot completed successfully.
func (t *PivotPagelinks) Done() bool {
	return lib.CheckExists(viper.GetString("pagelinks_pivoted"))
}

//Cleanup removes partial files on a failed pivot.
func (t *PivotPagelinks) Cleanup() error {
	return lib.CleanupFile(viper.GetString("pagelinks_pivoted"))
}

//Deps returns the dependencies of this task.
func (t *PivotPagelinks) Deps() []Task {
	return []Task{&SortPagelinks{}}
}
