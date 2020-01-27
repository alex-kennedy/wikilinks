package main

import (
	"flag"
	"log"

	"github.com/alex-kennedy/wikilinks"
	"github.com/alex-kennedy/wikilinks/prep/pipeline"
	"github.com/alex-kennedy/wikilinks/prep/tasks"
)

func main() {
	var configFileName = flag.String("config", "", "configuration file path")
	flag.Parse()

	wikilinks.InitialiseConfig(configFileName)

	err := pipeline.Run(&tasks.IndexPageDirect{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	// name := viper.GetString("page_direct_sorted")
	// index := viper.GetString("page_direct_index")

	// bSearcher, err := lib.NewBinarySearcher(name, index)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(bSearcher.Search("\"...\"\"Let_Me_Sing\"\"\""))
}
