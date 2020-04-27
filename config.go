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
	viper.Set("pagelinks", path.Join(dir, "pagelinks.csv"))
	viper.Set("pagelinks_resolved", path.Join(dir, "pagelinks_resolved.csv"))
	viper.Set("pagelinks_sorted", path.Join(dir, "pagelinks_sorted.csv"))
	viper.Set("pagelinks_pivoted", path.Join(dir, "pagelinks_pivoted.csv"))
	viper.Set("pagelinks_pivoted_index", path.Join(dir, "pagelinks_pivoted_index.csv"))
	viper.Set("backlinks", path.Join(dir, "backlinks.csv"))
	viper.Set("pagelinks_lgl", path.Join(dir, "pagelinks.lgl"))

	viper.Set("page_sql_gz", path.Join(dir, "page.sql.gz"))
	viper.Set("page_sql", path.Join(dir, "page.sql"))
	viper.Set("page_direct", path.Join(dir, "page_direct.csv"))
	viper.Set("page_redirect", path.Join(dir, "page_redirect.csv"))
	viper.Set("page_direct_sorted", path.Join(dir, "page_direct_sorted.csv"))
	viper.Set("page_direct_index", path.Join(dir, "page_direct_index.index"))
	viper.Set("page_redirect_sorted", path.Join(dir, "page_redirect_sorted.csv"))
	viper.Set("page_redirect_resolved", path.Join(dir, "page_redirect_resolved.csv"))
	viper.Set("page_merged", path.Join(dir, "page_merged.csv"))
	viper.Set("page_merged_index", path.Join(dir, "page_merged_index.index"))
	viper.Set("page_ids", path.Join(dir, "page_ids.txt"))

	viper.Set("redirect_sql_gz", path.Join(dir, "redirect.sql.gz"))
	viper.Set("redirect_sql", path.Join(dir, "redirect.sql"))
	viper.Set("redirect", path.Join(dir, "redirect.csv"))
	viper.Set("redirect_sorted", path.Join(dir, "redirect_sorted.csv"))
	viper.Set("redirect_index", path.Join(dir, "redirect_index.index"))

	viper.Set("pageranks", path.Join(dir, "pageranks.txt"))
}

//InitialiseConfig sets up the configuration with Viper.
func InitialiseConfig(configFileName *string) {
	// Overridden if config file was passed as a flag
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if *configFileName != "" {
		viper.SetConfigFile(*configFileName)
	}

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

	viper.SetDefault("bytes", 1024*1024*1024) // 1GB

	addFileNamesToConfig()
}
