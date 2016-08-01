package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func cleanAppName(t *testing.T) {
	assert.Equal(t, CleanAppName("/test.aayush.http"), "test.aayush.http", "They should be equal") // Dummy string
	assert.Equal(t, CleanAppName("/test/aayush/http"), "testaayushhttp", "They should be equal")
	assert.Equal(t, CleanAppName("/test.aayush.http/"), "test.aayush.http", "They should be equal")
}
