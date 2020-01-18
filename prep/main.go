package main

import (
	"flag"
	"log"

	"github.com/alex-kennedy/wikilinks"
	"github.com/spf13/viper"
)

func main() {
	var configFileName = flag.String("config", "", "configuration file path")
	flag.Parse()

	wikilinks.InitialiseConfig(configFileName)

	err := DownloadWikiFile("enwiki-20191201-redirect.sql.gz",
		"https://ftp.acc.umu.se/mirror/wikimedia.org/dumps/enwiki/20191201/",
		viper.GetString("redirect_sql_gz"))
	if err != nil {
		log.Fatalf("%s", err)
	}
}
