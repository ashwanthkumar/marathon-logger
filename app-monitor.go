package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ashwanthkumar/golang-utils/maps"
	marathon "github.com/gambol99/go-marathon"
)

type Task struct {
	App    string
	Labels map[string]string
	TaskID string
}

const LogEnabledLabel = "logs.enabled"

type AppMonitor struct {
	Client        marathon.Marathon
	RunWaitGroup  sync.WaitGroup
	CheckInterval time.Duration
	stopChannel   chan bool
	TasksChannel  chan Task // TaskInfo without CWD will be sent to this channel
}

func (a *AppMonitor) Start() {
	fmt.Println("Starting App Checker...")
	a.RunWaitGroup.Add(1)
	a.stopChannel = make(chan bool)
	a.TasksChannel = make(chan Task)
	go a.run()
	fmt.Println("App Checker Started.")
	fmt.Printf("App Checker - Checking the status of all the apps every %v\n", a.CheckInterval)
}

func (a *AppMonitor) Stop() {
	fmt.Println("Stopping App Checker...")
	close(a.stopChannel)
	a.RunWaitGroup.Done()
}

func (a *AppMonitor) run() {
	running := true
	for running {
		select {
		case <-time.After(a.CheckInterval):
			err := a.monitorApps()
			if err != nil {
				log.Fatalf("Unexpected error - %v\n", err)
			}
		case <-a.stopChannel:
			running = false
		}
		time.Sleep(1 * time.Second)
	}
}

func (a *AppMonitor) monitorApps() error {
	apps, err := a.Client.Applications(nil)
	if err != nil {
		return err
	}

	for _, app := range apps.Apps {
		isLogEnabled := maps.GetBoolean(app.Labels, LogEnabledLabel, false)
		if isLogEnabled {
			app, err := a.Client.Application(app.ID)
			if err != nil {
				return err
			}
			for _, task := range app.Tasks {
				taskInfo := Task{
					App:    app.ID,
					Labels: app.Labels,
					TaskID: task.ID,
				}
				a.TasksChannel <- taskInfo
			}
		}
	}

	return nil
}
