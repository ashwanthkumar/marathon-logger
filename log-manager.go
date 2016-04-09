package main

import (
	"fmt"
	"sync"
)

// TODO - Add support for using multiple logging backends
// LoggerToUse - Logging backend to use from app labes
// const LoggerToUse = "logger.backend"

// LogManager manages various logging backends
type LogManager struct {
	// Channel where we get tasks to start following
	Add chan TaskInfo
	// Channel where we get to idle tasks
	Remove  chan string
	Loggers map[string]Logger
	// DefaultLogger to use when none is specified in the app labels
	DefaultLogger string

	RunWaitGroup sync.WaitGroup
	stopChannel  chan bool
}

// Start - Start the LogManager
func (l *LogManager) Start() {
	fmt.Println("Starting Log Manager...")
	l.RunWaitGroup.Add(1)
	l.stopChannel = make(chan bool)
	go l.run()
	fmt.Println("Log Manager Started.")
}

// Stop - Stop the LogManager
func (l *LogManager) Stop() {
	fmt.Println("Stopping Log Manager...")
	close(l.stopChannel)

	l.RunWaitGroup.Done()
	fmt.Println("Log Manager Stoped.")
}

func (l *LogManager) run() {
	running := true
	for running {
		select {
		case addTaskInfo := <-l.Add:
			logger, present := l.Loggers[l.DefaultLogger]
			if !present {
				fmt.Printf("Requested logger %s is not found\n", l.DefaultLogger)
			} else {
				logger.AddTask(addTaskInfo)
			}
		case removeTaskInfo := <-l.Remove:
			logger, present := l.Loggers[l.DefaultLogger]
			if !present {
				fmt.Printf("Requested logger %s is not found\n", l.DefaultLogger)
			} else {
				logger.RemoveTask(removeTaskInfo)
			}
		case <-l.stopChannel:
			running = false
		}
	}
}
