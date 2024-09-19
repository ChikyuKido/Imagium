package jobs

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Job struct {
	Func    func()
	Rate    uint64 // in seconds
	lastRun time.Time
}

type JobHandler struct {
	jobs []Job
}

func (h *JobHandler) AddJob(job Job) {
	h.jobs = append(h.jobs, job)
}

func (h *JobHandler) Run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			for i := range h.jobs {
				job := &h.jobs[i]
				if now.Sub(job.lastRun).Seconds() >= float64(job.Rate) {
					logrus.Info("Run job")
					job.Func()
					job.lastRun = now
				}
			}
		}
	}
}
