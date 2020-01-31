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

	err := pipeline.Run(&tasks.ResolveRedirects{})
	if err != nil {
		log.Fatalf(err.Error())
	}
}
