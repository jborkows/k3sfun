package main

import (
	"fmt"
	"log"
	"net/http"
)

func debug(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}

func info(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func warning(format string, args ...interface{}) {
	log.Printf("[WARNING] "+format, args...)
}

func errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

func setAuthHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}

func expectJSON(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

func buildURL(config Config, path string) string {
	return config.APIURL + "/api/v1" + path
}

func buildTaskURL(config Config, taskID int, suffix string) string {
	if suffix != "" {
		return fmt.Sprintf("%s/api/v1/tasks/%d%s", config.APIURL, taskID, suffix)
	}
	return fmt.Sprintf("%s/api/v1/tasks/%d", config.APIURL, taskID)
}
