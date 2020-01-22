package pipeline

import (
	"errors"
	"fmt"

	"github.com/alex-kennedy/wikilinks/prep/tasks"
)

func TopoSortTasks(head tasks.Task) (*[]tasks.Task, error) {
	sorted := make([]tasks.Task, 0, 10)
	visited := make(map[string]bool)
	return visit(head, &sorted, visited)
}

func visit(node tasks.Task, sorted *[]tasks.Task, visited map[string]bool) (*[]tasks.Task, error) {
	name := fmt.Sprintf("%T", node)
	permanent, exists := visited[name]
	if exists && permanent {
		return sorted, nil
	}
	if exists && !permanent {
		return sorted, errors.New("tasks not a DAG")
	}
	visited[name] = false
	for _, dep := range node.Deps() {
		sorted, err := visit(dep, sorted, visited)
		if err != nil {
			return sorted, err
		}
	}
	visited[name] = true
	*sorted = append(*sorted, node)
	return sorted, nil
}
