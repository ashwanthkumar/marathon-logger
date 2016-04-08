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
# Created via marathon-alerts,
# PLEASE DON'T EDIT THIS FILE MANUALLY
# Name - {{ .App }}
# File - {{ .FileName }}
######################################

module(load="imfile")

input(type="imfile"
      File="{{ .CWD }}/{{ .FileName }}"
      Tag="{{ .App }}"
			statefile="{{ .TaskID }}"
      Severity="info")
`

type Rsyslog struct {
	ConfigLocation string
	RestartCommand string
}

func (r *Rsyslog) AddTask(taskInfo TaskInfo) error {
	fmt.Printf("[Rsyslog] Add task info for %v\n", taskInfo)
	template, err := r.render(taskInfo)
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
	return nil
}

func (r *Rsyslog) RemoveTask(taskId string) error {
	fmt.Printf("[Rsyslog] Remove task info for %v\n", taskId)
	return nil
}

func (r *Rsyslog) ExistingTasks() []string {
	return []string{}
}

func (r *Rsyslog) render(taskInfo TaskInfo) (string, error) {
	var configInBytes bytes.Buffer
	tmpl, err := template.New("").Parse(RsyslogTemplate)
	err = tmpl.Execute(&configInBytes, &taskInfo)
	return configInBytes.String(), err
}
