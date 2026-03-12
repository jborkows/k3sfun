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

func (a *AutoTransition) taskForBucket(bucketName BucketName) ([]Task, error) {
	bucketID, exists := a.BucketMapping[bucketName]
	if !exists {
		return nil, fmt.Errorf("bucket %q not found in bucket mapping", bucketName)
	}

	url := fmt.Sprintf("%s/api/v1/projects/%d/views/%d/tasks", a.APIURL, a.ProjectID, a.ViewID)

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

	type viewBucket struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Tasks []Task `json:"tasks"`
	}

	var viewBuckets []viewBucket
	if err := json.NewDecoder(resp.Body).Decode(&viewBuckets); err != nil {
		return nil, fmt.Errorf("failed to decode view tasks: %w", err)
	}

	for _, vb := range viewBuckets {
		if vb.ID == int(bucketID) {
			return vb.Tasks, nil
		}
	}

	return []Task{}, nil
}
