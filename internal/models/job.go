// internal/models/job.go
package models

import "time"

type Job struct {
	ID         string     `json:"job_id" gorm:"primaryKey"`
	Status     string     `json:"status" gorm:"default:QUEUED"`
	Progress   int        `json:"progress"`
	Processed  int        `json:"processed"`
	Total      int        `json:"total"`
	ResultPath string     `json:"result_path,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Cancelled  bool       `json:"cancelled"`
	CanceledAt *time.Time `json:"canceled_at,omitempty"`
}
