package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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

func (a *AutoTransition) get(url string) (*http.Response, error) {
	return a.doRequest("GET", url, nil)
}

func (a *AutoTransition) post(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	expectJSON(req)
	return a.doRequestFromReq(req)
}

func (a *AutoTransition) delete(url string) (*http.Response, error) {
	return a.doRequest("DELETE", url, nil)
}

func (a *AutoTransition) doRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	return a.doRequestFromReq(req)
}

func (a *AutoTransition) doRequestFromReq(req *http.Request) (*http.Response, error) {
	setAuthHeader(req, a.APIToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			errlog("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}
