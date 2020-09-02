package job

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

type (
	Task struct {
		TaskID uint64

		job *Job
	}

	TaskOption func(t *Task)

	CronTaskStatus string

	TaskResult struct {
		TaskID uint64         `json:"task_id"`
		Status CronTaskStatus `json:"status"`
		Job    *Job           `json:"job"`
		Logs   string         `json:"logs"`
		RunOn  string         `json:"run_on"`
	}
)

var (
	CronTaskStatusProcessing CronTaskStatus = "processing"
	CronTaskStatusSuccess    CronTaskStatus = "success"
	CronTaskStatusFailed     CronTaskStatus = "failed"
)

func NewTask(job *Job, ops ...TaskOption) *Task {
	now := time.Now()
	job.ExecutedAt = &now

	task := &Task{
		job: job,
	}
	for _, op := range ops {
		op(task)
	}

	if task.TaskID == 0 {
		id, _ := job.worker.taskIdGen.NextID()
		task.TaskID = id
	}

	return task
}

func (t *Task) SetStatus(status CronTaskStatus, logs string) error {
	if status == CronTaskStatusSuccess || status == CronTaskStatusFailed {
		now := time.Now()
		t.job.FinishedAt = &now
	}

	key := fmt.Sprintf("%s%s/%d", ResultKeyPrefix, t.job.ID, t.TaskID)

	payload := TaskResult{
		TaskID: t.TaskID,
		Job:    t.job,
		Status: status,
		Logs:   logs,
		RunOn:  t.job.runOn,
	}
	payloadBytes, _ := json.Marshal(&payload)

	session, err := concurrency.NewSession(t.job.Client.Client, concurrency.WithTTL(600))
	if err != nil {
		return err
	}

	_, err = t.job.Client.Put(context.Background(),
		key,
		string(payloadBytes),
		clientv3.WithLease(session.Lease()),
	)
	return err
}

func WithTaskID(taskId uint64) TaskOption {
	return func(t *Task) {
		t.TaskID = taskId
	}
}
