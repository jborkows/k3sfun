package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (a *AutoTransition) moveTaskToBucket(taskID int, targetBucketName BucketName) error {
	targetBucketID, exists := a.BucketMapping[targetBucketName]
	if !exists {
		return fmt.Errorf("bucket %q not found in bucket mapping", targetBucketName)
	}

	url := fmt.Sprintf("%s/api/v1/tasks/%d", a.APIURL, taskID)

	type updateRequest struct {
		BucketID int `json:"bucket_id"`
	}

	body, err := json.Marshal(updateRequest{BucketID: int(targetBucketID)})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	setAuthHeader(req, a.APIToken)
	expectJSON(req)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (a *AutoTransition) archiveOldTasks() {
	tasks, err := a.tasksToArchive()
	if err != nil {
		info("Error getting tasks to archive: %v", err)
		return
	}

	if len(tasks) == 0 {
		info("No tasks to archive")
		return
	}

	info("Found %d task(s) to archive:", len(tasks))
	for _, task := range tasks {
		info("  Would move task to archive: %s (ID: %d, DoneAt: %s)", task.Title, task.ID, task.DoneAt.Format(time.RFC3339))
	}
}
