package main

type BucketName string

const (
	awaitingBucket = BucketName("Awaiting")
	todoBucket     = BucketName("To-Do")
	doneBucket     = BucketName("Done")
	archiveBucket  = BucketName("Archive")
)

type BucketId int
type BucketMapping map[BucketName]BucketId

type Bucket struct {
	ID    BucketId   `json:"id"`
	Title BucketName `json:"title"`
}

type Task struct {
	ID           int          `json:"id"`
	Title        string       `json:"title"`
	Done         bool         `json:"done"`
	Labels       []string     `json:"labels"`
	StartDate    *string      `json:"start_date"`
	BucketID     *BucketId    `json:"bucket_id"`
	RelatedTasks RelatedTasks `json:"related_tasks"`
}

type RelatedTasks struct {
	Blocked  []BlockingTask `json:"blocked"`  // Tasks that block this task
	Blocking []BlockingTask `json:"blocking"` // Tasks that this task blocks
}

type BlockingTask struct {
	ID   int  `json:"id"`
	Done bool `json:"done"`
}

type TaskRelation struct {
	TaskID       int    `json:"task_id"`
	OtherTaskID  int    `json:"other_task_id"`
	RelationKind string `json:"relation_kind"`
}
