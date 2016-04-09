package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ashwanthkumar/golang-utils/maps"
	"github.com/ashwanthkumar/golang-utils/netutil"
	"github.com/ashwanthkumar/marathon-logger/mesos"
)

const LogFilesToMonitor = "logs.files"

type TaskInfo struct {
	App      string
	Labels   map[string]string
	TaskID   string
	Hostname string
	CWD      string // Current working directory of the task in the slave
	FileName string // Actual file name to that we need monitor for logs
}

// TaskManager - Enhances the Task with FileName and CWD info
// Message Flow: App Monitor -> Task Manager -> Log Manager
type TaskManager struct {
	InputTasksChannel         chan Task
	MaxTasksHeartBeatInterval time.Duration
	SlavePort                 int

	AddLogs    chan TaskInfo
	RemoveLogs chan string
	KnownTasks map[string]time.Time

	Client       mesos.Mesos
	RunWaitGroup sync.WaitGroup
	TasksMutex   sync.Mutex
	stopChannel  chan bool
}

// Start the TaskManager
func (t *TaskManager) Start() {
	fmt.Println("Starting Task Manager...")
	t.RunWaitGroup.Add(1)
	t.stopChannel = make(chan bool)
	t.AddLogs = make(chan TaskInfo)
	t.RemoveLogs = make(chan string)
	t.KnownTasks = make(map[string]time.Time)
	t.Client = &mesos.MesosClient{}
	go t.run()
	fmt.Println("Task Manager Started.")
	fmt.Printf("Task Manager - Task's MaxHeartBeatInterval is %v\n", t.MaxTasksHeartBeatInterval)
}

// Stop the TaskManager
func (t *TaskManager) Stop() {
	fmt.Println("Stopping Task Manager...")
	close(t.stopChannel)
	t.RunWaitGroup.Done()
}

func (t *TaskManager) run() {
	hostname, err := netutil.FullyQualifiedHostname()
	if err != nil {
		log.Fatalf("Error - %v\n", err)
	}
	running := true
	for running {
		select {
		case <-time.After(5 * time.Second):
			for task, lastHeartbeat := range t.KnownTasks {
				if time.Now().Sub(lastHeartbeat) > t.MaxTasksHeartBeatInterval {
					t.RemoveLogs <- task
					t.TasksMutex.Lock()
					delete(t.KnownTasks, task)
					t.TasksMutex.Unlock()
				}
			}
		case task := <-t.InputTasksChannel:
			if task.Hostname == hostname {
				// println("Got task for addition.. do what needs to be done")
				// fmt.Printf("%v\n", task)
				t.TasksMutex.Lock()
				_, present := t.KnownTasks[task.TaskID]
				if !present {
					fmt.Printf("TaskID %s is not monitored, sending it to LogManager", task.TaskID)
					slaveState, _ := t.Client.SlaveState(fmt.Sprintf("http://%s:%d/state.json", hostname, t.SlavePort))
					// fmt.Printf("%v\n", slaveState)
					executor := slaveState.FindExecutor(task.TaskID)
					if executor != nil {
						logFiles := strings.Split(maps.GetString(task.Labels, LogFilesToMonitor, "stdout"), ",")
						t.KnownTasks[task.TaskID] = time.Now()
						for _, file := range logFiles {
							taskInfo := TaskInfo{
								App:      task.App,
								Hostname: task.Hostname,
								Labels:   task.Labels,
								TaskID:   task.TaskID,
								CWD:      executor.Directory,
								FileName: file,
							}
							// fmt.Printf("%v\n", taskInfo)
							t.AddLogs <- taskInfo
						}
					} else {
						fmt.Printf("[WARN] Couldn't find the executor that spun up the task %s", task.TaskID)
					}
				} else {
					// Already present - update the clock
					t.KnownTasks[task.TaskID] = time.Now()
				}
				t.TasksMutex.Unlock()
			}
			// 1. Check if the task is running on our machine
			// 2. Check if we already know the task
			// 3. If yes, just update the KnownTasks map
			// 4. Else, construct a TaskInfo object and send it to Log Manager and update the KnownTasks map
		case <-t.stopChannel:
			running = false
		}
	}
}
