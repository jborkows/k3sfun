package main

import (
	"time"
)

func main() {
	config := loadConfigWithToken()

	info("Auto-transition running with check interval: %v", config.CheckInterval)

	// Run immediately
	processTasks(config)

	// Then run periodically
	ticker := time.NewTicker(config.CheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		processTasks(config)
	}
}
