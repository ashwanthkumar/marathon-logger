package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"text/template"
)

// CommonPrefixToConfigFiles - Common Prefix To Config Files
const CommonPrefixToConfigFiles = "marathon-logger"

// RsyslogTemplate - Go Template for Rsyslog configuration
// TODO - Make this configurable
const RsyslogTemplate = `
######################################
# Created via marathon-logger,
# PLEASE DON'T EDIT THIS FILE MANUALLY
# Name - {{ .App }}
# File - {{ .FileName }}
######################################

module(load="imfile")

input(type="imfile"
			File="{{ .WorkDir }}/{{ .FileName }}"
			Tag="{{ .CleanAppName }}	{{.TaskID}}"
			statefile="{{ .TaskID }}"
      Severity="info")
`

// Rsyslog backend implementation
type Rsyslog struct {
	WorkDir string
	ConfigLocation string
	RestartCommand string
}

// AddTask - Adds a task definition file to FS
func (r *Rsyslog) AddTask(taskInfo TaskInfo) error {
	fmt.Printf("[Rsyslog] Add task info for %v\n", taskInfo)
	template, err := r.render(taskInfo)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return err
	}

	// We create symlink of the tasks' sandbox dir so that, 
	// we can reduce the length of the file path that we provide as 
	// part of configs in rsyslog's imfile directive.
	// When the full file path > 200, rsyslog crashes with "buffer overflow" (evil smile)
	err = os.Symlink(taskInfo.CWD, fmt.Sprintf("%s/%s", r.WorkDir, taskInfo.TaskID))
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return err
	}

	configFileLocation := fmt.Sprintf("%s/%s-%s.conf", r.ConfigLocation, CommonPrefixToConfigFiles, taskInfo.TaskID)
	err = ioutil.WriteFile(configFileLocation, []byte(template), 0644)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return err
	}

	err = exec.Command("/bin/sh", "-c", r.RestartCommand).Run()
	return err
}

// RemoveTask - Remove a task definition from the FS
func (r *Rsyslog) RemoveTask(taskId string) error {
	fmt.Printf("[Rsyslog] Remove task info for %v\n", taskId)
	return nil
}

// TODO - Integrate it with LogManager
func (r *Rsyslog) ExistingTasks() []string {
	return []string{}
}

func (r *Rsyslog) render(taskInfo TaskInfo) (string, error) {
	var configInBytes bytes.Buffer
	tmpl, err := template.New("").Parse(RsyslogTemplate)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&configInBytes, &taskInfo)
	return configInBytes.String(), err
}
