package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/alex-kennedy/wikilinks"
	"github.com/alex-kennedy/wikilinks/prep/pipeline"
	"github.com/alex-kennedy/wikilinks/prep/status"
	"github.com/alex-kennedy/wikilinks/prep/tasks"
)

func main() {
	var configFileName = flag.String("config", "", "configuration file path")
	var portNumber = flag.String("port", "80", "Port for status server")
	flag.Parse()

	wikilinks.InitialiseConfig(configFileName)

	rootTask := &tasks.SaveBacklinks{}
	taskPipeline, err := pipeline.NewPipeline(rootTask)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Main pipeline
	go taskPipeline.Run()

	// Status webserver
	statusSite := status.NewStatusSite(taskPipeline)
	http.HandleFunc("/", statusSite.RenderPage)
	http.ListenAndServe(":"+*portNumber, nil)
}
