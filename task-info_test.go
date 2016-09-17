package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanAppName(t *testing.T) {
	var testCaseOne, testCaseTwo, testCaseThree TaskInfo
	testCaseOne.App = "/test.aayush.http"
	assert.Equal(t, testCaseOne.CleanAppName(), "test.aayush.http", "They should be equal") // Dummy string
	testCaseTwo.App = "/test/aayush/http"
	assert.Equal(t, testCaseTwo.CleanAppName(), "test.aayush.http", "They should be equal")
	testCaseThree.App = "/test.aayush.http/"
	assert.Equal(t, testCaseThree.CleanAppName(), "test.aayush.http.", "They should be equal")
}
