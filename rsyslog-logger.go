package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
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
# File - {{ .FileNames }}
######################################

module(load="imfile")
{{range $fileName := .FileNames}}
input(type="imfile"
	File="{{ $.WorkDir }}/{{$fileName}}"
	Tag="{{$.CleanAppName}}	{{$.TaskID}}	{{$fileName}}"
	Severity="info")
{{end}}
`

// Rsyslog backend implementation
type Rsyslog struct {
	WorkDir              string
	SyslogConfigLocation string
	RestartCommand       string
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
	err = os.Symlink(taskInfo.CWD, r.symlink(taskInfo.TaskID))
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return err
	}

	configFileLocation := fmt.Sprintf("%s/%s-%s.conf", r.SyslogConfigLocation, CommonPrefixToConfigFiles, taskInfo.TaskID)
	err = ioutil.WriteFile(configFileLocation, []byte(template), 0644)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return err
	}

	err = exec.Command("/bin/sh", "-c", r.RestartCommand).Run()
	return err
}

// RemoveTask removes a task definition from the FS
func (r *Rsyslog) RemoveTask(taskId string) error {
	fmt.Printf("[Rsyslog] Remove task info for %v\n", taskId)
	// remove the rsyslog conf file
	rsyslogConfig := fmt.Sprintf("%s/%s-%s.conf", r.SyslogConfigLocation, CommonPrefixToConfigFiles, taskId)
	err := os.Remove(rsyslogConfig)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return err
	}

	// remove the symlink created in the workding dir
	symlinkFile := r.symlink(taskId)
	err = os.Remove(symlinkFile)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return err
	}

	err = exec.Command("/bin/sh", "-c", r.RestartCommand).Run()
	return err
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
	taskInfo.WorkDir = r.symlink(taskInfo.TaskID)
	err = tmpl.Execute(&configInBytes, &taskInfo)
	return configInBytes.String(), err
}

func (r *Rsyslog) symlink(taskId string) string {
	return fmt.Sprintf("%s/%s", r.WorkDir, taskId)
}
