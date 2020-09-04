package job

import "runtime"

// 单次任务
type OnceJob struct {
	Job

	TaskID uint64 `json:"task_id"`
}

func (o *OnceJob) RunWithRecovery(taskOptions ...TaskOption) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			o.logger.Warnf("panic running job: %v\n%s", r, buf)
		}
	}()
	_ = o.Run(taskOptions...)
}
