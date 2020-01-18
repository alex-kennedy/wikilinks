package wikilinks

import (
	"fmt"
	"path"

	"github.com/spf13/viper"
)

func addFileNamesToConfig() {
	viper.SetDefault("root_dir", "./data/")

	rootDir := viper.GetString("root_dir")

	viper.Set("pagelinks_sql_gz", path.Join(rootDir, "pagelinks.sql.gz"))
	viper.Set("pagelinks_sql", path.Join(rootDir, "pagelinks.sql"))

	viper.Set("page_sql_gz", path.Join(rootDir, "page.sql.gz"))
	viper.Set("page_sql", path.Join(rootDir, "page.sql"))

	viper.Set("redirect_sql_gz", path.Join(rootDir, "redirect.sql.gz"))
	viper.Set("redirect_sql", path.Join(rootDir, "redirect.sql"))
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
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	}

	addFileNamesToConfig()
}
