package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type JobRecord struct {
	ID         string     `gorm:"primaryKey" json:"job_id"`
	Status     string     `json:"status"`   
	Progress   int        `json:"progress"`
	Processed  int64      `json:"processed"`
	Total      int64      `json:"total"`
	ResultPath string     `json:"result_path"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Cancelled  bool       `json:"cancelled"`
	CancelAt   *time.Time `json:"cancel_at,omitempty"`
}

type JobRepository interface {
	Create(ctx context.Context, job *JobRecord) error
	UpdateStatus(ctx context.Context, id string, status string) error
	GetByID(ctx context.Context, id string) (*JobRecord, error)
	MarkCancelled(ctx context.Context, id string) error
	UpdateJob(ctx context.Context, job *JobRecord) error
	UpdateTotal(ctx context.Context, jobID string, total int64) error
	UpdateProgress(ctx context.Context, jobID string, processed int, progress int) error
	IncrementProcessed(ctx context.Context, id string, inc int64, progress float64) error
}

type jobRepo struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepo{db: db}
}

func (r *jobRepo) Create(ctx context.Context, job *JobRecord) error {
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO job_records (
			id, status, progress, processed, total, result_path, created_at, updated_at, cancelled, cancel_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			status = VALUES(status),
			updated_at = VALUES(updated_at),
			progress = VALUES(progress),
			processed = VALUES(processed),
			total = VALUES(total),
			result_path = VALUES(result_path)
	`, job.ID, job.Status, job.Progress, job.Processed, job.Total, job.ResultPath, job.CreatedAt, job.UpdatedAt, job.Cancelled, job.CancelAt).Error
}

func (r *jobRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&JobRecord{}).Where("id = ?", id).Update("status", status).Error
}

func (r *jobRepo) GetByID(ctx context.Context, id string) (*JobRecord, error) {
	var job JobRecord
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&job).Error; err != nil {
		return nil, err
	}
	if job.Status == "FINISHED" {
		job.ResultPath = "/public/downloads/" + job.ID + ".csv"
	}
	return &job, nil
}

func (r *jobRepo) MarkCancelled(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&JobRecord{}).Where("id = ?", id).Updates(map[string]interface{}{
		"cancelled":  true,
		"status":     "CANCELED",
		"updated_at": time.Now(),
	}).Error
}

func (r *jobRepo) UpdateJob(ctx context.Context, job *JobRecord) error {
	// Pastikan status dan timestamp benar
	job.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).
		Model(&JobRecord{}).
		Where("id = ?", job.ID).
		Updates(map[string]interface{}{
			"progress":   job.Progress,
			"processed":  job.Processed,
			"total":      job.Total,
			"status":     job.Status,
			"updated_at": job.UpdatedAt,
		}).Error
}

func (r *jobRepo) UpdateTotal(ctx context.Context, jobID string, total int64) error {
	return r.db.WithContext(ctx).
		Model(&JobRecord{}).
		Where("id = ?", jobID).
		Update("total", total).Error
}

func (r *jobRepo) UpdateProgress(
	ctx context.Context,
	jobID string,
	processed int,
	progress int,
) error {
	return r.db.WithContext(ctx).
		Model(&JobRecord{}).
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"processed":  processed,
			"progress":   progress,
			"updated_at": time.Now(),
		}).Error
}

func (r *jobRepo) IncrementProcessed(ctx context.Context, id string, inc int64, progress float64) error {
	return r.db.WithContext(ctx).
		Model(&JobRecord{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"processed":  gorm.Expr("processed + ?", inc),
			"progress":   progress,
			"updated_at": time.Now(),
		}).Error
}
