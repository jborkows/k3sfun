package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	labelStateReady      = "state:ready"
	labelStateBlocked    = "state:blocked"
	labelStateScheduled  = "state:scheduled"
	labelStateInProgress = "state:in-progress"
	labelStateOnHalt     = "state:on-halt"
	labelStateCompleted  = "state:completed"

	prefixEarliestOn = "earliest-on:"
)

type Task struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Done   bool     `json:"done"`
	Labels []string `json:"labels"`
}

type TaskRelation struct {
	TaskID       int    `json:"task_id"`
	OtherTaskID  int    `json:"other_task_id"`
	RelationKind string `json:"relation_kind"`
}

func processTasks(config Config) {
	info("Running task checks...")
	processBlockedTasks(config)
	processScheduledTasks(config)
	info("Task checks complete")
}

func processBlockedTasks(config Config) {
	info("Checking blocked tasks...")

	tasks, err := getTasksWithLabel(config, labelStateBlocked)
	if err != nil {
		errorf("Failed to get blocked tasks: %v", err)
		return
	}

	for _, task := range tasks {
		info("Checking blocked task: %s (ID: %d)", task.Title, task.ID)

		hasBlockers, err := hasIncompleteBlockers(config, task.ID)
		if err != nil {
			errorf("Failed to check blockers for task %d: %v", task.ID, err)
			continue
		}

		if hasBlockers {
			info("Still has incomplete blockers")
			continue
		}

		info("No incomplete blockers, transitioning to ready")

		newLabels := updateLabels(task.Labels, labelStateBlocked, labelStateReady)
		if err := updateTaskLabels(config, task.ID, newLabels); err != nil {
			errorf("Failed to update task: %v", err)
			continue
		}
		info("Successfully transitioned to ready")
	}
}

func processScheduledTasks(config Config) {
	info("Checking scheduled tasks...")

	today := time.Now().Format("2006-01-02")

	tasks, err := getTasksWithLabel(config, labelStateScheduled)
	if err != nil {
		errorf("Failed to get scheduled tasks: %v", err)
		return
	}

	for _, task := range tasks {
		scheduledDate := extractEarliestOnDate(task.Labels)
		if scheduledDate == "" {
			continue
		}

		info("Task: %s (ID: %d) - scheduled for: %s", task.Title, task.ID, scheduledDate)

		if today < scheduledDate {
			info("Still waiting for %s (today: %s)", scheduledDate, today)
			continue
		}

		info("Scheduled date reached, transitioning to ready")

		newLabels := updateScheduledLabels(task.Labels, labelStateReady)
		if err := updateTaskLabels(config, task.ID, newLabels); err != nil {
			errorf("Failed to update task: %v", err)
			continue
		}
		info("Successfully transitioned to ready")
	}
}

func getTasksWithLabel(config Config, label string) ([]Task, error) {
	url := buildURL(config, "/tasks?filter_labels="+label)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	setAuthHeader(req, config.APIToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			warning("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Tasks []Task `json:"tasks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		// Try decoding as array directly
		var tasks []Task
		if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
			return []Task{}, nil
		}
		return tasks, nil
	}

	return result.Tasks, nil
}

func hasIncompleteBlockers(config Config, taskID int) (bool, error) {
	url := buildTaskURL(config, taskID, "/relations")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	setAuthHeader(req, config.APIToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			warning("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var relations []TaskRelation
	if err := json.NewDecoder(resp.Body).Decode(&relations); err != nil {
		return false, fmt.Errorf("failed to decode relations: %w", err)
	}

	for _, rel := range relations {
		if rel.RelationKind == "blocks" {
			done, err := isTaskDone(config, rel.OtherTaskID)
			if err != nil {
				return false, fmt.Errorf("failed to check task %d: %w", rel.OtherTaskID, err)
			}
			if !done {
				return true, nil
			}
		}
	}

	return false, nil
}

func isTaskDone(config Config, taskID int) (bool, error) {
	url := buildTaskURL(config, taskID, "")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	setAuthHeader(req, config.APIToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return false, err
	}

	return task.Done, nil
}

func updateTaskLabels(config Config, taskID int, labels []string) error {
	url := fmt.Sprintf("%s/api/v1/tasks/%d", config.APIURL, taskID)

	payload := map[string][]string{"labels": labels}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	setAuthHeader(req, config.APIToken)
	expectJSON(req)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func updateLabels(labels []string, removeLabel, addLabel string) []string {
	newLabels := []string{}

	// Remove old label and add new one
	for _, label := range labels {
		if label != removeLabel {
			newLabels = append(newLabels, label)
		}
	}

	// Add new label if not already present
	hasNewLabel := false
	for _, label := range newLabels {
		if label == addLabel {
			hasNewLabel = true
			break
		}
	}

	if !hasNewLabel {
		newLabels = append(newLabels, addLabel)
	}

	return newLabels
}

func updateScheduledLabels(labels []string, addLabel string) []string {
	newLabels := []string{}

	// Remove state:scheduled and earliest-on:* labels
	for _, label := range labels {
		if label != labelStateScheduled && !strings.HasPrefix(label, prefixEarliestOn) {
			newLabels = append(newLabels, label)
		}
	}

	// Add new label if not already present
	hasNewLabel := false
	for _, label := range newLabels {
		if label == addLabel {
			hasNewLabel = true
			break
		}
	}

	if !hasNewLabel {
		newLabels = append(newLabels, addLabel)
	}

	return newLabels
}

func extractEarliestOnDate(labels []string) string {
	for _, label := range labels {
		if strings.HasPrefix(label, prefixEarliestOn) {
			return strings.TrimPrefix(label, prefixEarliestOn)
		}
	}
	return ""
}
