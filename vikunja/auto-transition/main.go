package main

import (
	"log"
	"time"
)

func main() {

	config := loadConfigWithToken()

	info("Auto-transition running with check interval: %v", config.CheckInterval)
	info("Project ID: %d, View ID: %d", config.ProjectID, config.ViewID)

	at, err := NewAutoTransition(config)
	if err != nil {
		log.Fatalf("Cannot initialize processor %v", err)
	}

	at.Run()

	ticker := time.NewTicker(config.CheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		at.Run()
	}
}
