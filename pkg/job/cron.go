// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package job

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/douyu/jupiter/pkg/util/xstring"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type (
	// JobWrapper ...
	JobWrapper = cron.JobWrapper
	// EntryID ...
	EntryID = cron.EntryID
	// Schedule ...
	Schedule = cron.Schedule
	// Job ...
	CronJob = cron.Job
	//NamedJob ..
	NamedJob interface {
		Run() error
	}
)

// FuncJob ...
type FuncJob func() error

// Run ...
func (f FuncJob) Run() error { return f() }

// Name ...
func (f FuncJob) Name() string { return xstring.FunctionName(f) }

type wrappedLogger struct {
	*xlog.Logger
}

// Info logs routine messages about cron's operation.
func (wl *wrappedLogger) Info(msg string, keysAndValues ...interface{}) {
	wl.Sugar().Infow("cron "+msg, keysAndValues...)
}

// Error logs an error condition.
func (wl *wrappedLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	wl.Sugar().Errorw("cron "+msg, append(keysAndValues, "err", err)...)
}

// Cron ...
type Cron struct {
	*Worker
	*cron.Cron
	entries map[string]EntryID
}

func newCron(config *Worker) *Cron {
	c := &Cron{
		Worker: config,
		Cron: cron.New(
			cron.WithLogger(&wrappedLogger{config.logger}),
			cron.WithChain(config.wrappers...),
		),
	}
	return c
}

// Schedule ...
func (c *Cron) Schedule(schedule Schedule, job NamedJob) EntryID {
	if c.ImmediatelyRun {
		schedule = &immediatelyScheduler{
			Schedule: schedule,
		}
	}
	innnerJob := &wrappedJob{
		NamedJob: job,
		logger:   c.Worker.logger,
	}

	return c.Cron.Schedule(schedule, innnerJob)
}

// AddJob ...
func (c *Cron) AddJob(spec string, cmd NamedJob) (EntryID, error) {
	schedule, err := c.Worker.parser.Parse(spec)
	if err != nil {
		return 0, err
	}
	return c.Schedule(schedule, cmd), nil
}

// AddFunc ...
func (c *Cron) AddFunc(spec string, cmd func() error) (EntryID, error) {
	return c.AddJob(spec, FuncJob(cmd))
}

// Remove an entry from being run in the future.
func (c *Cron) Remove(id EntryID) {
	c.Cron.Remove(id)
}

// Run ...
func (c *Cron) Run() {
	c.Worker.logger.Info("run Worker", xlog.Int("number of scheduled jobs", len(c.Cron.Entries())))
	c.Cron.Start()
}

// Stop ...
func (c *Cron) Stop() error {
	_ = c.Cron.Stop()
	return nil
}

type immediatelyScheduler struct {
	Schedule
	initOnce uint32
}

// Next ...
func (is *immediatelyScheduler) Next(curr time.Time) (next time.Time) {
	if atomic.CompareAndSwapUint32(&is.initOnce, 0, 1) {
		return curr
	}

	return is.Schedule.Next(curr)
}

type wrappedJob struct {
	NamedJob
	logger *xlog.Logger
}

// Run ...
func (wj wrappedJob) Run() {
	_ = wj.run()
}

func (wj wrappedJob) run() (err error) {
	var fields = []xlog.Field{}
	var beg = time.Now()
	defer func() {
		if rec := recover(); rec != nil {
			switch rec := rec.(type) {
			case error:
				err = rec
			default:
				err = fmt.Errorf("%v", rec)
			}

			stack := make([]byte, 4096)
			length := runtime.Stack(stack, true)
			fields = append(fields, zap.ByteString("stack", stack[:length]))
		}
		if err != nil {
			fields = append(fields, xlog.String("err", err.Error()), xlog.Duration("cost", time.Since(beg)))
			wj.logger.Error("Worker", fields...)
		}
	}()

	return wj.NamedJob.Run()
}
