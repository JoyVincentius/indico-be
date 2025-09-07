// internal/job/job_queue.go (versi perbaikan)
package job

import (
	"context"
	"encoding/json"
	"fmt"
	"indico-be/internal/repository"
	"sync"
	"time"

	"github.com/google/uuid"
)

// JobQueue is the façade used by HTTP handlers.
type JobQueue struct {
	queue      chan *Job
	workers    []*Worker
	workerPool *WorkerPool
	mu         sync.Mutex
	jobRepo    repository.JobRepository
}

// NewJobQueue creates a channel + workers.
func NewJobQueue(pool *WorkerPool) *JobQueue {
	q := &JobQueue{
		queue:      make(chan *Job, 100),
		workerPool: pool,
	}
	// attach workers
	for i := 0; i < pool.Count; i++ {
		w := NewWorker(i+1, pool.Service, q.queue)
		w.Start()
		q.workers = append(q.workers, w)
	}
	return q
}

// generateJobID returns a UUID‑v4 string.
func generateJobID() string {
	return uuid.NewString()
}

func (jq *JobQueue) SetRepository(repo repository.JobRepository) {
	jq.jobRepo = repo
}

// Enqueue creates a Job record and pushes to channel.
func (q *JobQueue) Enqueue(from, to string) (string, error) {
	// --------- 1️⃣ Parse tanggal ----------
	fromT, err := time.Parse("2006-01-02", from)
	if err != nil {
		return "", fmt.Errorf("invalid from date: %w", err)
	}
	toT, err := time.Parse("2006-01-02", to)
	if err != nil {
		return "", fmt.Errorf("invalid to date: %w", err)
	}

	// --------- 2️⃣ Buat objek Job ----------
	j := &Job{
		ID:        generateJobID(),
		From:      fromT,
		To:        toT,
		CreatedAt: time.Now(),
	}

	// 🔥 HAPUS: TIDAK PERLU SAVE ke DB di sini!
	// JANGAN panggil q.jobRepo.Create

	// --------- 3️⃣ Kirim ke channel ---------- (HANYA ini yang diperlukan!)
	q.queue <- j
	return j.ID, nil
}

// Cancel a running job.
func (q *JobQueue) Cancel(jobID string) error {
	// Mark cancelled in DB – workers will notice via the job record.
	return q.jobRepo.MarkCancelled(context.Background(), jobID)
}

// Close shuts down the channel (used on graceful shutdown).
func (q *JobQueue) Close() {
	close(q.queue)
}

// Helper to return JSON status for API.
func (q *JobQueue) Status(jobID string) ([]byte, error) {
	rec, err := q.jobRepo.GetByID(context.Background(), jobID)
	if err != nil {
		return nil, err
	}
	return json.Marshal(rec)
}
