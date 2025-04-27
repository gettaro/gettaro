package scheduler

import (
	"context"
	"time"
)

type Job interface {
	Run(ctx context.Context) error
}

type Scheduler struct {
	job      Job
	interval time.Duration
}

func NewScheduler(job Job, interval time.Duration) *Scheduler {
	return &Scheduler{
		job:      job,
		interval: interval,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run immediately
	if err := s.job.Run(ctx); err != nil {
		// Log error but continue
	}

	// Then run on schedule
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.job.Run(ctx); err != nil {
				// Log error but continue
			}
		}
	}
}
