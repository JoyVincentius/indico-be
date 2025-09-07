package handler

import (
	"net/http"

	"indico-be/internal/job"
	"indico-be/internal/repository"

	"github.com/gin-gonic/gin"
)

type settlementReq struct {
	From string `json:"from" binding:"required"` // "2025-01-01"
	To   string `json:"to" binding:"required"`
}

// RegisterJobRoutes sets up /jobs endpoints.
func RegisterJobRoutes(r *gin.Engine, q *job.JobQueue, repo repository.JobRepository) {
	jobs := r.Group("/jobs")
	{
		jobs.POST("/settlement", submitJob(q))
		jobs.GET("/:id", getJobStatus(q))
		jobs.POST("/:id/cancel", cancelJob(q))
		jobs.GET("/downloads/:filename", serveCSV())
	}
}

func submitJob(q *job.JobQueue) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req settlementReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ✅ Asumsi: Enqueue hanya mengirim tanggal, bukan objek job
		jobID, err := q.Enqueue(req.From, req.To)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ✅ Tidak ada CreateJob di handler!
		c.JSON(http.StatusAccepted, gin.H{
			"job_id": jobID,
			"status": "QUEUED",
		})
	}
}

func getJobStatus(q *job.JobQueue) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		b, err := q.Status(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
			return
		}
		c.Data(http.StatusOK, "application/json", b)
	}
}

func cancelJob(q *job.JobQueue) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := q.Cancel(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"job_id": id, "status": "CANCELLING"})
	}
}

// Simple static file serving for generated CSVs.
func serveCSV() gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := "public/downloads/" + filename
		c.File(filePath)
	}
}
