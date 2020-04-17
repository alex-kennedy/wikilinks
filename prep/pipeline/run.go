package pipeline

import (
	"fmt"
	"log"

	"github.com/alex-kennedy/wikilinks/prep/tasks"
)

//TaskStatus captures the states which tasks may be in.
type TaskStatus string

//Constants pertaining to task status.
const (
	DoneStatus       TaskStatus = "done"
	InProgressStatus TaskStatus = "inProgress"
	NotDoneStatus    TaskStatus = "notDone"
)

//Pipeline holds information about task state and implements Run to run the
//tasks.
type Pipeline struct {
	RootTask        tasks.Task
	PipelineRunning bool
	PipelineError   bool
	TaskStatuses    map[string]TaskStatus
	SortedTasks     *[]tasks.Task
}

//Run runs the RootTask after any required dependencies.
func (p *Pipeline) Run() error {
	for _, t := range *p.SortedTasks {
		if !t.Done() {
			p.TaskStatuses[fmt.Sprintf("%T", t)] = InProgressStatus
			err := t.Run()
			if err != nil {
				t.Cleanup()
				p.PipelineError = true
				p.PipelineRunning = false
				return err
			}
			p.TaskStatuses[fmt.Sprintf("%T", t)] = DoneStatus
		} else {
			log.Printf("%T done, skipping", t)
			p.TaskStatuses[fmt.Sprintf("%T", t)] = DoneStatus
		}
	}
	p.PipelineRunning = false
	return nil
}

//NewPipeline initiates the Pipeline object with a RootTask to run.
func NewPipeline(task tasks.Task) (*Pipeline, error) {
	sortedTasks, err := TopoSortTasks(task)
	if err != nil {
		return nil, err
	}

	taskStatues := make(map[string]TaskStatus)
	log.Println("Tasks sorted")
	for i, t := range *sortedTasks {
		taskStatues[fmt.Sprintf("%T", t)] = NotDoneStatus
		log.Printf("%d. %T\n", i+1, t)
	}

	return &Pipeline{task, true, false, taskStatues, sortedTasks}, nil
}
