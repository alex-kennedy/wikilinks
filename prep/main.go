package main

import (
	"flag"
	"log"

	"github.com/alex-kennedy/wikilinks"
	"github.com/alex-kennedy/wikilinks/prep/tasks"
)

func main() {
	var configFileName = flag.String("config", "", "configuration file path")
	flag.Parse()

	wikilinks.InitialiseConfig(configFileName)

	task := tasks.ExtractRedirect{}
	err := task.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
