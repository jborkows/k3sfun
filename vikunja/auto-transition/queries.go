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

func (a *AutoTransition) tasksToArchive() ([]Task, error) {
	tasks, err := a.taskForBucket(doneBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks from done bucket: %w", err)
	}

	var result []Task
	todayStart := time.Now().Truncate(24 * time.Hour)

	for _, task := range tasks {
		if !task.Done {
			continue
		}
		if task.DoneAt == nil {
			continue
		}
		if task.DoneAt.Before(todayStart) {
			result = append(result, task)
		}
	}

	return result, nil
}

func (a *AutoTransition) blockedTasksInTodo() ([]Task, error) {
	todoTasks, err := a.taskForBucket(todoBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks from todo bucket: %w", err)
	}

	activeBuckets := []BucketName{todoBucket, doingBucket, pendingBucket}
	activeTaskIDs := make(map[int]bool)

	for _, bucketName := range activeBuckets {
		tasks, err := a.taskForBucket(bucketName)
		if err != nil {
			return nil, fmt.Errorf("failed to get tasks from %s bucket: %w", bucketName, err)
		}
		for _, task := range tasks {
			activeTaskIDs[task.ID] = true
		}
	}

	var result []Task
	for _, task := range todoTasks {
		for _, blocker := range task.RelatedTasks.Blocked {
			if activeTaskIDs[blocker.ID] {
				result = append(result, task)
				break
			}
		}
	}

	return result, nil
}

func (a *AutoTransition) unblockedTasksInAwaiting() ([]Task, error) {
	awaitingTasks, err := a.taskForBucket(awaitingBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks from awaiting bucket: %w", err)
	}

	doneTasks, err := a.taskForBucket(doneBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks from done bucket: %w", err)
	}

	doneTaskIDs := make(map[int]bool)
	for _, task := range doneTasks {
		doneTaskIDs[task.ID] = true
	}

	var result []Task
	for _, task := range awaitingTasks {
		if len(task.RelatedTasks.Blocked) == 0 {
			continue
		}

		allBlockersDone := true
		for _, blocker := range task.RelatedTasks.Blocked {
			if !doneTaskIDs[blocker.ID] {
				allBlockersDone = false
				break
			}
		}

		if allBlockersDone {
			result = append(result, task)
		}
	}

	return result, nil
}

func (a *AutoTransition) tasksToDeleteFromArchive() ([]Task, error) {
	tasks, err := a.taskForBucket(archiveBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks from archive bucket: %w", err)
	}

	threshold := time.Now().AddDate(0, 0, -7)

	var result []Task
	for _, task := range tasks {
		if !task.Done {
			continue
		}
		if task.DoneAt == nil {
			continue
		}
		if task.DoneAt.Before(threshold) {
			result = append(result, task)
		}
	}

	return result, nil
}
