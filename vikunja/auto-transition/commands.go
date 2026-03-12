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

	// Use view-specific bucket endpoint for Kanban
	url := fmt.Sprintf("%s/api/v1/projects/%d/views/%d/buckets/%d/tasks",
		a.APIURL, a.ProjectID, a.ViewID, targetBucketID)

	type taskBucketRequest struct {
		TaskID int `json:"task_id"`
	}

	body, err := json.Marshal(taskBucketRequest{TaskID: taskID})
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
		info("Moving task to archive: %s (ID: %d, DoneAt: %s)", task.Title, task.ID, task.DoneAt.Format(time.RFC3339))
		if err := a.moveTaskToBucket(task.ID, archiveBucket); err != nil {
			info("  ERROR: Failed to move task: %v", err)
		} else {
			info("  Successfully moved to archive bucket")
		}
	}
}
