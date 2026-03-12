package main

import (
	"bytes"
	"os"
	"strconv"
	"time"
)

const (
	defaultAPIURL        = "http://localhost:3456"
	defaultCheckInterval = 1 * 60 * 60
	defaultProjectID     = 2
	defaultViewID        = 8

	envAPIURL        = "VIKUNJA_API_URL"
	envAPIToken      = "VIKUNJA_API_TOKEN"
	envCheckInterval = "CHECK_INTERVAL"
	envProjectID     = "PROJECT_ID"
	envViewID        = "VIEW_ID"
)

type autoTransitionConfig struct {
	APIURL        string
	APIToken      string
	CheckInterval time.Duration
	ProjectID     int
	ViewID        int
}

type Config = autoTransitionConfig

func loadConfig() autoTransitionConfig {
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

	projectID := defaultProjectID
	if val := os.Getenv(envProjectID); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			projectID = parsed
		}
	}

	viewID := defaultViewID
	if val := os.Getenv(envViewID); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			viewID = parsed
		}
	}

	apiToken := os.Getenv(envAPIToken)
	if apiToken == "" {
		if data, err := os.ReadFile("/secrets/API_TOKEN"); err == nil {
			apiToken = string(bytes.TrimSpace(data))
		}
	}

	return autoTransitionConfig{
		APIURL:        apiURL,
		APIToken:      apiToken,
		CheckInterval: time.Duration(interval) * time.Second,
		ProjectID:     projectID,
		ViewID:        viewID,
	}
}

func loadConfigWithToken() autoTransitionConfig {
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
