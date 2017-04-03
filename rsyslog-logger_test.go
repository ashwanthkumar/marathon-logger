package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderRsyslogTemplate(t *testing.T) {

	hostname, err := os.Hostname()
	var rsyslog = Rsyslog{
		WorkDir: "/foo/bar",
	}
	label := map[string]string{
		"logs.enabled": "true",
	}
	taskInfo := TaskInfo{
		App:      "/test.aayush.http",
		Hostname: hostname,
		Labels:   label,
		TaskID:   "abcdefghij",
		CWD:      "/foo/bar",
		FileNames: []string{"test_file_name.txt"},
	}

	expected := `
######################################
# Created via marathon-logger,
# PLEASE DON'T EDIT THIS FILE MANUALLY
# Name - /test.aayush.http
# File - [test_file_name.txt]
######################################

module(load="imfile")

input(type="imfile"
	File="/foo/bar/abcdefghij/test_file_name.txt"
	Tag="test.aayush.http	abcdefghij	test_file_name.txt"
	Severity="info")

`
	template, err := rsyslog.render(taskInfo)
	assert.NoError(t, err)
	assert.Equal(t, expected, template)
}

func TestRenderRsyslogTemplateForMultipleFiles(t *testing.T) {

	hostname, err := os.Hostname()
	var rsyslog = Rsyslog{
		WorkDir: "/foo/bar",
	}
	label := map[string]string{
		"logs.enabled": "true",
	}
	taskInfo := TaskInfo{
		App:      "/test.aayush.http",
		Hostname: hostname,
		Labels:   label,
		TaskID:   "abcdefghij",
		CWD:      "/foo/bar",
		FileNames: []string{"test_file_name1.txt","test_file_name2.txt"},
	}

	expected := `
######################################
# Created via marathon-logger,
# PLEASE DON'T EDIT THIS FILE MANUALLY
# Name - /test.aayush.http
# File - [test_file_name1.txt test_file_name2.txt]
######################################

module(load="imfile")

input(type="imfile"
	File="/foo/bar/abcdefghij/test_file_name1.txt"
	Tag="test.aayush.http	abcdefghij	test_file_name1.txt"
	Severity="info")

input(type="imfile"
	File="/foo/bar/abcdefghij/test_file_name2.txt"
	Tag="test.aayush.http	abcdefghij	test_file_name2.txt"
	Severity="info")

`
	template, err := rsyslog.render(taskInfo)
	assert.NoError(t, err)
	assert.Equal(t, expected, template)
}
