package valueobject

import (
	"fmt"
	"time"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/repository/model"
)

// TaskStageName is the name of the task stage.
type TaskStageName string

const (
	TaskStageApproval     TaskStageName = "approval" // authenticity review
	TaskStageRewrite      TaskStageName = "rewrite"  // rewrite it as an colloquial podcast article
	TaskStageClassify     TaskStageName = "classify" // classify the topic
	TaskStageStylize      TaskStageName = "stylize"  // stylize the article
	TaskStageMerge        TaskStageName = "merge"    // merge the article
	TaskStageScripted     TaskStageName = "scripted" // scripted
	TaskStageTextToSpeech TaskStageName = "tts"      // text to speech
)

// StagePriority is the priority of the task stage.
var StagePriority = []TaskStageName{
	TaskStageApproval,
	TaskStageRewrite,
	TaskStageClassify,
	TaskStageScripted,
}

// StageStatus is the status of the task stage.
type StageStatus string

const (
	StageStatusProcessing StageStatus = "processing"
	StageStatusCompleted  StageStatus = "completed"
	StageStatusFailed     StageStatus = "failed"
)

// TaskStageExtra represents extra information for the task stage.
type TaskStageExtra struct {
	NewsIds  []uint
	StageIds []uint
}

// TaskStage represents a task stage.
type TaskStage struct {
	Id        uint
	BatchNo   string
	Stage     TaskStageName
	Status    StageStatus
	Prompt    string
	Input     string
	Output    string
	Reason    string
	Audio     *PodcastAudio
	TaskAi    *TaskAi
	Extra     *TaskStageExtra
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewTaskStage creates a new task stage.
func NewTaskStage(stage TaskStageName, prompt string, ai *TaskAi) *TaskStage {
	return &TaskStage{
		Stage:     stage,
		Status:    StageStatusProcessing,
		Prompt:    prompt,
		TaskAi:    ai,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewTaskStageFromModel creates a task stage from a model.PodcastTask.
func NewTaskStageFromModel(m *model.PodcastTask) (*TaskStage, error) {
	return &TaskStage{
		Id:        m.ID,
		BatchNo:   m.BatchNo,
		Stage:     TaskStageName(m.Stage),
		Status:    StageStatus(m.Status),
		Prompt:    m.Prompt,
		Input:     m.Input,
		Output:    m.Output,
		Reason:    m.Reason,
		Audio:     gokit.UnmarshalSafe[*PodcastAudio](m.Audio),
		TaskAi:    gokit.UnmarshalSafe[*TaskAi](m.TaskAi),
		Extra:     gokit.UnmarshalSafe[*TaskStageExtra](m.Extra),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

// Fail marks the task stage as failed.
func (s *TaskStage) Fail(reason string) {
	s.Status = StageStatusFailed
	s.Reason = reason
	s.UpdatedAt = time.Now()
}

// SetStatus sets the status of the task stage.
func (s *TaskStage) SetStatus(status StageStatus) {
	s.Status = status
	s.UpdatedAt = time.Now()
}

// SetOutput sets the output of the task stage.
func (s *TaskStage) SetOutput(output string) {
	if s.Status == StageStatusProcessing {
		s.Status = StageStatusCompleted
	}

	s.Output = output
	s.UpdatedAt = time.Now()
}

// BuildPrompt builds the prompt for the task stage.
func (s *TaskStage) BuildPrompt() string {
	prompt := s.Prompt

	if s.Input != "" {
		prompt = fmt.Sprintf("%s\n%s", prompt, s.Input)
	}

	return prompt
}
