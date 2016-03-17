package main

type Logger interface {
	// Add a configuration for a mesos task in the Logger backend
	AddTask(taskInfo TaskInfo) error
	// Remove a configuration for a mesos task in the Logger backend
	RemoveTask(taskInfo TaskInfo) error
	// Get the list of all known tasks already configured in the logger backend
	// This is used only to get the state back when we die and come back
	ExistingTasks() []string
}
