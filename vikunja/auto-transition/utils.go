package main

import (
	"log"
	"net/http"
)

func info(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func warning(format string, args ...interface{}) {
	log.Printf("[WARNING] "+format, args...)
}

func errlog(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

func setAuthHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}

func expectJSON(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}
