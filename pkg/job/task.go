package job

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/douyu/jupiter/pkg/xlog"
)

type (
	Task struct {
		TaskID uint64

		job        *Job
		executedAt time.Time
		finishedAt *time.Time
	}

	TaskOption func(t *Task)

	CronTaskStatus string

	TaskResult struct {
		TaskID     uint64         `json:"task_id"`
		Status     CronTaskStatus `json:"status"`
		Job        *Job           `json:"job"`
		Logs       string         `json:"logs"`
		RunOn      string         `json:"run_on"`
		ExecutedAt time.Time      `json:"executed_at"`
		FinishedAt *time.Time     `json:"finished_at"`
	}
)

var (
	CronTaskStatusProcessing CronTaskStatus = "processing"
	CronTaskStatusSuccess    CronTaskStatus = "success"
	CronTaskStatusFailed     CronTaskStatus = "failed"
	CronTaskStatusTimeout    CronTaskStatus = "timeout"
)

func NewTask(job *Job, ops ...TaskOption) *Task {
	task := &Task{
		job:        job,
		executedAt: time.Now(),
	}
	for _, op := range ops {
		op(task)
	}
	if task.TaskID == 0 {
		id, _ := job.Worker.taskIdGen.NextID()
		task.TaskID = id
	}

	return task
}

func (t *Task) SetStatus(status CronTaskStatus, logs string) error {
	if status == CronTaskStatusSuccess || status == CronTaskStatusFailed || status == CronTaskStatusTimeout {
		now := time.Now()
		t.finishedAt = &now
	}

	payload := TaskResult{
		TaskID:     t.TaskID,
		Job:        t.job,
		Status:     status,
		Logs:       logs,
		RunOn:      t.job.HostName,
		ExecutedAt: t.executedAt,
		FinishedAt: t.finishedAt,
	}
	payloadBytes, _ := json.Marshal(&payload)

	_, err := t.job.Client.Put(context.Background(),
		t.Key(),
		string(payloadBytes),
	)
	return err
}

func (t *Task) Key() string {
	return fmt.Sprintf("%s%s/%d", ResultKeyPrefix, t.job.ID, t.TaskID)
}

func (t *Task) Stop() {
	_, err := t.job.Client.Delete(context.Background(), t.Key())
	if err != nil {
		t.job.logger.Error("delete task result failed", xlog.FieldErr(err))
	}
}

func WithTaskID(taskId uint64) TaskOption {
	return func(t *Task) {
		t.TaskID = taskId
	}
}
