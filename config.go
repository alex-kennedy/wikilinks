package wikilinks

import (
	"fmt"
	"path"
	"strings"

	"github.com/spf13/viper"
)

func addFileNamesToConfig() {
	viper.SetDefault("root_dir", "./data/")

	dir := path.Join(viper.GetString("root_dir"), viper.GetString("date"))

	viper.Set("pagelinks_sql_gz", path.Join(dir, "pagelinks.sql.gz"))
	viper.Set("pagelinks_sql", path.Join(dir, "pagelinks.sql"))

	viper.Set("page_sql_gz", path.Join(dir, "page.sql.gz"))
	viper.Set("page_sql", path.Join(dir, "page.sql"))

	viper.Set("redirect_sql_gz", path.Join(dir, "redirect.sql.gz"))
	viper.Set("redirect_sql", path.Join(dir, "redirect.sql"))
	viper.Set("redirect", path.Join(dir, "redirect.csv"))
}

//InitialiseConfig sets up the configuration with Viper.
func InitialiseConfig(configFileName *string) {
	if *configFileName != "" {
		viper.SetConfigFile(*configFileName)
	}

	// Overridden if config file was passed as a flag
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Config file not found"))
		} else {
			panic(fmt.Errorf("Fatal error in config file: %s", err))
		}
	}

	if viper.GetString("date") == "" {
		panic("No date in config")
	}

	siteURL := viper.GetString("site_url")
	if siteURL == "" {
		panic("No site_url in config")
	}
	if !strings.HasSuffix(siteURL, "/") {
		siteURL = siteURL + "/"
	}

	addFileNamesToConfig()
}
