package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func buckets(a autoTransitionConfig) (BucketMapping, error) {

	url := fmt.Sprintf("%s/api/v1/projects/%d/views/%d/buckets", a.APIURL, a.ProjectID, a.ViewID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	setAuthHeader(req, a.APIToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var buckets []Bucket
	if err := json.NewDecoder(resp.Body).Decode(&buckets); err != nil {
		return nil, fmt.Errorf("failed to decode buckets: %w", err)
	}

	var bucketCache BucketMapping = make(BucketMapping)
	for _, b := range buckets {
		bucketCache[BucketName(b.Title)] = BucketId(b.ID)
	}

	return bucketCache, nil
}
