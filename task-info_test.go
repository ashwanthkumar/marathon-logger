package main

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestCleanAppName(t *testing.T) {
	var testCaseOne, testCaseTwo, testCaseThree TaskInfo
	testCaseOne.App = "/test.aayush.http"
	assert.Equal(t, testCaseOne.CleanAppName(), "test.aayush.http", "They should be equal") // Dummy string
	testCaseTwo.App = "/test/aayush/http"
	assert.Equal(t, testCaseTwo.CleanAppName(), "test.aayush.http", "They should be equal")
	testCaseThree.App = "/test.aayush.http"
	assert.Equal(t, testCaseThree.CleanAppName(), "test.aayush.http", "They should be equal")
}

func TestRenderRsyslogTemplate(t *testing.T) {

	hostname, err := os.Hostname()
	var rsyslog Rsyslog
	label := map[string]string{
		"logs.enabled": "true",
	}
	taskInfo := TaskInfo{
		App:      "/test.aayush.http",
		Hostname: hostname,
		Labels:   label,
		TaskID:   randSeq(10),
		CWD:      "/foo/bar",
		FileName: "test_file_name.txt",
	}

	template, err := rsyslog.render(taskInfo)
	assert.NoError(t, err)
	fmt.Printf("%s", template)
}
