package main

import (
	"time"
)

type BucketName string

const (
	awaitingBucket = BucketName("Awaiting")
	todoBucket     = BucketName("To-Do")
	doneBucket     = BucketName("Done")
	doingBucket    = BucketName("Doing")
	pendingBucket  = BucketName("Pending")
)

type BucketId int
type TaskId int
type BucketMapping map[BucketName]BucketId
type TaskSet map[TaskId]bool

func (s TaskSet) Add(id TaskId) {
	s[id] = true
}

func (s TaskSet) Contains(id TaskId) bool {
	return s[id]
}

type Bucket struct {
	ID    BucketId   `json:"id"`
	Title BucketName `json:"title"`
}

type Task struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Done      bool       `json:"done"`
	DoneAt    *time.Time `json:"done_at"`
	Labels    []string   `json:"labels"`
	StartDate *string    `json:"start_date"`
	BucketID  *BucketId  `json:"bucket_id"`
}

type TaskWithRelations struct {
	Task
	RelatedTasks RelatedTasks `json:"related_tasks"`
}

type RelatedTasks struct {
	Blocked  []BlockingTask `json:"blocked"`  // Tasks that this task is blocked by (incoming blocking relations)
	Blocking []BlockingTask `json:"blocking"` // Tasks that this task blocks (outgoing blocking relations)
}

type BlockingTask struct {
	ID   TaskId `json:"id"`
	Done bool   `json:"done"`
}

type TaskRelation struct {
	TaskID       int    `json:"task_id"`
	OtherTaskID  int    `json:"other_task_id"`
	RelationKind string `json:"relation_kind"`
}
