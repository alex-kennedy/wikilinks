package main

import (
	"flag"
	"fmt"

	"github.com/alex-kennedy/wikilinks/prep/lib"

	"github.com/alex-kennedy/wikilinks"
)

func main() {
	var configFileName = flag.String("config", "", "configuration file path")
	flag.Parse()

	wikilinks.InitialiseConfig(configFileName)

	// err := pipeline.Run(&tasks.ExtractPagelinks{})
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }
	err := lib.IndexFile("data/20191201/test", "data/20191201/testindex.csv")
	fmt.Println(err)
}
