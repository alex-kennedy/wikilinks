package main

//Task is not done
type Task interface {
	run()
	cleanup()
	deps()
}
