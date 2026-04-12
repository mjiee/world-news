package valueobject

// TaskResult
type TaskResult string

const (
	TaskResultCompleted TaskResult = "completed"
	TaskResultFailed    TaskResult = "failed"
)

func (t TaskResult) IsCompleted() bool {
	return t == TaskResultCompleted
}

func (t TaskResult) IsFailed() bool {
	return t == TaskResultFailed
}
