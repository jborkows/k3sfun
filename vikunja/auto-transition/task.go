package main

import (
	"fmt"
)

type AutoTransition struct {
	autoTransitionConfig
	BucketMapping
}

func NewAutoTransition(config autoTransitionConfig) (*AutoTransition, error) {
	buckets, err := buckets(config)
	if err != nil {
		return nil, fmt.Errorf("While initializing %w", err)
	}

	return &AutoTransition{
		autoTransitionConfig: config,
		BucketMapping:        buckets,
	}, nil
}

func (a *AutoTransition) Run() {
	info("Running task checks...")
	info("%v", a.BucketMapping)

	// Test taskForBucket with archiveBucket
	tasks, err := a.taskForBucket(archiveBucket)
	if err != nil {
		info("Error fetching tasks for archiveBucket: %v", err)
	} else {
		info("Tasks in archiveBucket:")
		for _, t := range tasks {
			info("  Task ID: %d, Title: %s", t.ID, t.Title)
		}
	}

	info("Task checks complete")
}
