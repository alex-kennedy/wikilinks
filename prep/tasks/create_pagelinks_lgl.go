package tasks

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/lib"
	"github.com/spf13/viper"
)

//CreatePagelinksLGL converts the pagelinks_pivoted.csv file to a recognised
//.lgl file format. Requires the pagelinks to be loaded into memory so any
//bi-directional links can be eliminated.
type CreatePagelinksLGL struct{}

//Run converts the pagelinks file.
func (t *CreatePagelinksLGL) Run() error {
	log.Println("Converting pagelinks to LGL...")
	pagelinksPivotedName := viper.GetString("pagelinks_pivoted")
	pageIDsName := viper.GetString("page_ids")
	outName := viper.GetString("pagelinks_lgl")
	return lib.CreatePagelinksLGL(pagelinksPivotedName, pageIDsName, outName)
}

//Done checks if the conversion completed successfully.
func (t *CreatePagelinksLGL) Done() bool {
	return lib.CheckExists(viper.GetString("pagelinks_lgl"))
}

//Cleanup removes partial files on a failed conversion.
func (t *CreatePagelinksLGL) Cleanup() error {
	return lib.CleanupFile(viper.GetString("pagelinks_lgl"))
}

//Deps returns the dependencies of this task.
func (t *CreatePagelinksLGL) Deps() []Task {
	return []Task{&PivotPagelinks{}}
}
