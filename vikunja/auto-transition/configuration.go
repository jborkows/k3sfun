package main

import (
	"fmt"
	"os"
	"time"
)

const (
	defaultAPIURL        = "http://localhost:3456"
	defaultCheckInterval = 300 // 5 minutes

	envAPIURL        = "VIKUNJA_API_URL"
	envAPIToken      = "VIKUNJA_API_TOKEN"
	envCheckInterval = "CHECK_INTERVAL"
)

type Config struct {
	APIURL        string
	APIToken      string
	CheckInterval time.Duration
}

func loadConfig() Config {
	apiURL := os.Getenv(envAPIURL)
	if apiURL == "" {
		apiURL = defaultAPIURL
	}

	interval := defaultCheckInterval
	if val := os.Getenv(envCheckInterval); val != "" {
		if parsed, err := time.ParseDuration(val + "s"); err == nil {
			interval = int(parsed.Seconds())
		}
	}

	return Config{
		APIURL:        apiURL,
		APIToken:      os.Getenv(envAPIToken),
		CheckInterval: time.Duration(interval) * time.Second,
	}
}

func loadConfigWithToken() Config {
	config := loadConfig()

	info("Vikunja Auto-Transition Sidecar starting...")
	info("API URL: %s", config.APIURL)
	info("Check interval: %v", config.CheckInterval)

	if config.APIToken == "" {
		warning("VIKUNJA_API_TOKEN not set. Auto-transition disabled.")
		info("Waiting for API token configuration...")

		for {
			time.Sleep(60 * time.Second)
			token := os.Getenv(envAPIToken)
			if token != "" {
				config.APIToken = token
				info("API token detected, auto-transition enabled.")
				break
			}
		}
	}

	return config
}

func validateConfig(config Config) error {
	if config.APIURL == "" {
		return fmt.Errorf("API URL is required")
	}
	return nil
}
