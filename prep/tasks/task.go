package tasks

//Task is not done
type Task interface {
	Run() error
	Cleanup() error
	Deps() []Task
	Done() bool
}
