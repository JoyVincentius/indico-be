package job

import (
	"context"
	"log"
	"time"

	"indico-be/internal/repository"
	"indico-be/internal/service"
)

type Worker struct {
	id      int
	svc     *service.SettlementService
	jobChan <-chan *Job
}

func NewWorker(id int, svc *service.SettlementService, jobChan <-chan *Job) *Worker {
	return &Worker{id: id, svc: svc, jobChan: jobChan}
}

func (w *Worker) Start() {
	go func() {
		for job := range w.jobChan {
			ctx, cancel := context.WithCancel(context.Background())
			job.Cancel = cancel

			log.Printf("[worker %d] started job %s", w.id, job.ID)

			rec := &repository.JobRecord{
				ID:         job.ID,
				Status:     "RUNNING",
				Progress:   0,
				Processed:  0,
				Total:      0,
				ResultPath: "",
				CreatedAt:  job.CreatedAt,
				UpdatedAt:  time.Now(),
				Cancelled:  false,
				CancelAt:   nil,
			}

			if err := w.svc.JobRepo.Create(context.Background(), rec); err != nil {
				log.Printf("[worker %d] gagal membuat job di DB: %v", w.id, err)
				continue 
			}

			if err := w.svc.RunJob(ctx, job.ID, job.From.Format("2006-01-02"), job.To.Format("2006-01-02")); err != nil {
				log.Printf("[worker %d] job %s gagal: %v", w.id, job.ID, err)

				_ = w.svc.JobRepo.UpdateStatus(context.Background(), job.ID, "FAILED")
			} else {
				log.Printf("[worker %d] job %s berhasil", w.id, job.ID)

				_ = w.svc.JobRepo.UpdateStatus(context.Background(), job.ID, "FINISHED")
			}
		}
	}()
}
