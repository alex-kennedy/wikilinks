package pipeline

import (
	"log"

	"github.com/alex-kennedy/wikilinks/prep/tasks"
)

func Run(task tasks.Task) error {
	sorted, err := TopoSortTasks(task)
	if err != nil {
		return err
	}

	log.Println("Tasks sorted")
	for i, t := range *sorted {
		log.Printf("%d. %T\n", i+1, t)
	}

	for _, t := range *sorted {
		if !t.Done() {
			err := t.Run()
			if err != nil {
				t.Cleanup()
				return err
			}
		} else {
			log.Printf("%T done, skipping", t)
		}
	}
	return nil
}
