package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	marathon "github.com/gambol99/go-marathon"
	flag "github.com/spf13/pflag"
)

var marathonURI string
var mesosSlavePort int
var appCheckInterval time.Duration
var taskMaxHeartBeatInterval time.Duration
var rsyslogConfigurationDir string
var rsyslogRestartCommand string

var appMonitor AppMonitor
var taskManager TaskManager
var logManager LogManager

func main() {
	os.Args[0] = "marathon-logger"
	flag.Parse()

	client, err := marathonClient(marathonURI)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	appMonitor = AppMonitor{
		Client:        client,
		CheckInterval: appCheckInterval,
	}
	appMonitor.Start()

	taskManager = TaskManager{
		SlavePort:                 mesosSlavePort,
		InputTasksChannel:         appMonitor.TasksChannel,
		MaxTasksHeartBeatInterval: taskMaxHeartBeatInterval,
	}
	taskManager.Start()

	loggers := make(map[string]Logger)
	loggers["rsyslog"] = &Rsyslog{
		WorkDir: workDir,
		SyslogConfigLocation: rsyslogConfigurationDir,
		RestartCommand: rsyslogRestartCommand,
	}
	logManager = LogManager{
		Add:           taskManager.AddLogs,
		Remove:        taskManager.RemoveLogs,
		Loggers:       loggers,
		DefaultLogger: "rsyslog",
	}
	logManager.Start()

	appMonitor.RunWaitGroup.Wait()
}

func init() {
	flag.StringVar(&marathonURI, "uri", "", "Marathon URI to connect")
	flag.IntVar(&mesosSlavePort, "slave-port", 5051, "Mesos slave port")
	flag.DurationVar(&appCheckInterval, "app-check-interval", 30*time.Second, "Frequency at which we check for new tasks")
	flag.DurationVar(&taskMaxHeartBeatInterval, "task-max-heart-beat-interval", 30*time.Minute, "Max heartbeat interval after which the task is considered dead and logger is removed")
	flag.StringVar(&workDir, "work-dir", "/tmp/", "Location on the Filesystem where we create symlinks to app location base dir. This is needed to ensure the file paths don't become too long and crashes rsyslog")
	flag.StringVar(&rsyslogConfigurationDir, "rsyslog-configuration-dir", "/etc/rsyslog.d", "Location on the Filesystem where the rsyslog configurations needs to be written")
	flag.StringVar(&rsyslogRestartCommand, "rsyslog-restart-cmd", "service rsyslog restart", "Restart command for rsyslog backend")
}

func marathonClient(uri string) (marathon.Marathon, error) {
	config := marathon.NewDefaultConfig()
	config.URL = uri
	config.HTTPClient = &http.Client{
		Timeout: (30 * time.Second),
	}

	return marathon.NewClient(config)
}
