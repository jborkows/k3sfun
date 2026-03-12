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
	info("Task checks complete")
}
