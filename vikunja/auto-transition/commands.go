package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

func (a *AutoTransition) moveTaskToBucket(taskID int, targetBucketName BucketName) error {
	targetBucketID, exists := a.BucketMapping[targetBucketName]
	if !exists {
		return fmt.Errorf("bucket %q not found in bucket mapping", targetBucketName)
	}

	url := fmt.Sprintf("%s/api/v1/projects/%d/views/%d/buckets/%d/tasks",
		a.APIURL, a.ProjectID, a.ViewID, targetBucketID)

	type taskBucketRequest struct {
		TaskID int `json:"task_id"`
	}

	body, err := json.Marshal(taskBucketRequest{TaskID: taskID})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	_, err = a.post(url, bytes.NewReader(body))
	return err
}

func (a *AutoTransition) deleteTask(taskID int) error {
	url := fmt.Sprintf("%s/api/v1/tasks/%d", a.APIURL, taskID)
	_, err := a.delete(url)
	return err
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

func (a *AutoTransition) moveBlockedTasksToAwaiting() {
	blockedTasks, err := a.blockedTasksInTodo()
	if err != nil {
		info("Error getting blocked tasks in todo: %v", err)
		return
	}

	if len(blockedTasks) == 0 {
		info("No blocked tasks in todo bucket")
		return
	}

	info("Found %d blocked task(s) in todo bucket:", len(blockedTasks))
	for _, task := range blockedTasks {
		info("Moving blocked task to awaiting: %s (ID: %d)", task.Title, task.ID)
		if err := a.moveTaskToBucket(task.ID, awaitingBucket); err != nil {
			info("  ERROR: Failed to move task: %v", err)
		} else {
			info("  Successfully moved to awaiting bucket")
		}
	}
}

func (a *AutoTransition) moveUnblockedTasksToTodo() {
	unblockedTasks, err := a.unblockedTasksInAwaiting()
	if err != nil {
		info("Error getting unblocked tasks in awaiting: %v", err)
		return
	}

	if len(unblockedTasks) == 0 {
		info("No unblocked tasks in awaiting bucket")
		return
	}

	info("Found %d unblocked task(s) in awaiting bucket:", len(unblockedTasks))
	for _, task := range unblockedTasks {
		info("Moving unblocked task to todo: %s (ID: %d)", task.Title, task.ID)
		if err := a.moveTaskToBucket(task.ID, todoBucket); err != nil {
			info("  ERROR: Failed to move task: %v", err)
		} else {
			info("  Successfully moved to todo bucket")
		}
	}
}

func (a *AutoTransition) deleteOldArchivedTasks() {
	tasks, err := a.tasksToDeleteFromArchive()
	if err != nil {
		info("Error getting tasks to delete from archive: %v", err)
		return
	}

	if len(tasks) == 0 {
		info("No archived tasks older than 7 days to delete")
		return
	}

	info("Found %d archived task(s) older than 7 days to delete:", len(tasks))
	for _, task := range tasks {
		info("  Deleting: %s (ID: %d, DoneAt: %s)", task.Title, task.ID, task.DoneAt.Format(time.RFC3339))
		if err := a.deleteTask(task.ID); err != nil {
			info("    ERROR: Failed to delete task: %v", err)
		} else {
			info("    Successfully deleted task")
		}
	}
}
