package pipeline

import (
	"github.com/alex-kennedy/wikilinks/prep/tasks"
)

func Run(task tasks.Task) error {
	sorted, err := TopoSortTasks(task)
	if err != nil {
		return err
	}
	for _, t := range *sorted {
		if !t.Done() {
			err := t.Run()
			if err != nil {
				t.Cleanup()
				return err
			}
		}
	}
	return nil
}
