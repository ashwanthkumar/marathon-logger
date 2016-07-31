package main

import "strings"

// CleanAppName drops first `/` character from app name
func CleanAppName(appName string) string {
	return strings.Replace(appName, "/", "", 1)
}
