package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func buckets(a autoTransitionConfig) (BucketMapping, error) {
	url := fmt.Sprintf("%s/api/v1/projects/%d/views/%d/buckets", a.APIURL, a.ProjectID, a.ViewID)

	temp := &AutoTransition{
		autoTransitionConfig: a,
		BucketMapping:        make(BucketMapping),
	}

	resp, err := temp.get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get buckets: %w", err)
	}
	defer cleanup(resp)

	var buckets []Bucket
	if err := json.NewDecoder(resp.Body).Decode(&buckets); err != nil {
		return nil, fmt.Errorf("failed to decode buckets: %w", err)
	}

	bucketCache := make(BucketMapping)
	for _, b := range buckets {
		bucketCache[BucketName(b.Title)] = BucketId(b.ID)
	}

	return bucketCache, nil
}

func (a *AutoTransition) taskDetailFor(taskID int) (*TaskWithRelations, error) {
	url := fmt.Sprintf("%s/api/v1/tasks/%d", a.APIURL, taskID)

	resp, err := a.get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	defer cleanup(resp)

	var task TaskWithRelations
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode task: %w", err)
	}

	return &task, nil
}

func (a *AutoTransition) taskForBucket(bucketName BucketName) ([]Task, error) {
	bucketID, exists := a.BucketMapping[bucketName]
	if !exists {
		return nil, fmt.Errorf("bucket %q not found in bucket mapping", bucketName)
	}

	url := fmt.Sprintf("%s/api/v1/projects/%d/views/%d/tasks", a.APIURL, a.ProjectID, a.ViewID)

	resp, err := a.get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer cleanup(resp)

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

	activeTaskIDs := a.activeTasks()

	var result []Task
	for _, task := range todoTasks {
		fullTask, err := a.taskDetailFor(task.ID)
		if err != nil {
			warning("Failed to get task details for ID %d: %v", task.ID, err)
			continue
		}

		if isBlockedByActive(*fullTask, activeTaskIDs) {
			result = append(result, fullTask.Task)
		}
	}

	return result, nil
}

func isBlockedByActive(task TaskWithRelations, activeTaskIDs TaskSet) bool {
	for _, blocker := range task.RelatedTasks.Blocked {
		if activeTaskIDs.Contains(blocker.ID) {
			return true
		}
	}
	return false
}

func (a *AutoTransition) activeTasks() TaskSet {
	activeBuckets := []BucketName{todoBucket, doingBucket, pendingBucket}
	activeTaskIDs := make(TaskSet)

	for _, bucketName := range activeBuckets {
		tasks, _ := a.taskForBucket(bucketName)
		for _, task := range tasks {
			activeTaskIDs.Add(TaskId(task.ID))
		}
	}

	return activeTaskIDs
}

func (a *AutoTransition) unblockedTasksInAwaiting() ([]Task, error) {
	awaitingTasks, err := a.taskForBucket(awaitingBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks from awaiting bucket: %w", err)
	}

	doneTaskIDs := a.doneTasks()

	var result []Task
	for _, task := range awaitingTasks {
		fullTask, err := a.taskDetailFor(task.ID)
		if err != nil {
			warning("Failed to get task details for ID %d: %v", task.ID, err)
			continue
		}

		blockerIDs := getBlockerIDs(*fullTask)
		if len(blockerIDs) == 0 {
			continue
		}

		if allBlockersDone(blockerIDs, doneTaskIDs) {
			result = append(result, fullTask.Task)
		}
	}

	return result, nil
}

func getBlockerIDs(task TaskWithRelations) []TaskId {
	var blockerIDs []TaskId
	for _, b := range task.RelatedTasks.Blocked {
		blockerIDs = append(blockerIDs, b.ID)
	}
	return blockerIDs
}

func allBlockersDone(blockerIDs []TaskId, doneTaskIDs TaskSet) bool {
	for _, blockerID := range blockerIDs {
		if !doneTaskIDs.Contains(blockerID) {
			return false
		}
	}
	return true
}

func (a *AutoTransition) doneTasks() TaskSet {
	doneTasks, _ := a.taskForBucket(doneBucket)
	doneTaskIDs := make(TaskSet)
	for _, task := range doneTasks {
		doneTaskIDs.Add(TaskId(task.ID))
	}
	return doneTaskIDs
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
