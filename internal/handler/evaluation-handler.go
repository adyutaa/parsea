package handler

import (
	"net/http"
	"github.com/adyutaa/parsea/internal/service"

	"github.com/gin-gonic/gin"
)

type EvaluationHandler struct {
	service *service.EvaluationService
}

func NewEvaluationHandler(service *service.EvaluationService) *EvaluationHandler {
	return &EvaluationHandler{service: service}
}

// EvaluateRequest represents the request body for evaluation
type EvaluateRequest struct {
	CVID     string `json:"cv_id" binding:"required"`
	ReportID string `json:"report_id" binding:"required"`
	JobTitle string `json:"job_title" binding:"required"`
}

// Evaluate creates a new evaluation job
func (h *EvaluationHandler) Evaluate(c *gin.Context) {
	var req EvaluateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// Start evaluation
	jobID, err := h.service.StartEvaluation(req.CVID, req.ReportID, req.JobTitle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start evaluation: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     jobID,
		"status": "queued",
	})
}

// GetResult retrieves the result of an evaluation job
func (h *EvaluationHandler) GetResult(c *gin.Context) {
	jobID := c.Query("id")

	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Job ID is required as query parameter: ?id=your-job-id",
		})
		return
	}

	job, err := h.service.GetJobStatus(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Job not found",
			"job_id": jobID,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": job.ID,
		"status": job.Status,
		"result": job.Result,
		"created_at": job.CreatedAt,
		"updated_at": job.UpdatedAt,
	})
}

// GetQueueStatus returns current queue information
func (h *EvaluationHandler) GetQueueStatus(c *gin.Context) {
	length, err := h.service.GetQueueLength()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get queue status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"queue_length": length,
		"status":       "active",
	})
}
