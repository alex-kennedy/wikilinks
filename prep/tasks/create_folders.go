package tasks

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

//CreateFolders ensures the data directory exists.
type CreateFolders struct{}

//Run creates directory if it doesn't exist
func (t *CreateFolders) Run() error {
	path := path.Join(viper.GetString("root_dir"), viper.GetString("date"))
	return os.MkdirAll(path, 0700)
}

//Done checks that the folder exists
func (t *CreateFolders) Done() bool {
	path := path.Join(viper.GetString("root_dir"), viper.GetString("date"))
	info, err := os.Stat(path)
	return !os.IsNotExist(err) && info.IsDir()
}

//Cleanup in this case does nothing, we don't mind leftover directories.
func (t *CreateFolders) Cleanup() error {
	return nil
}

//Deps returns the dependencies of this task.
func (t *CreateFolders) Deps() []Task {
	return []Task{}
}
