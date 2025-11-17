package handler

import (
	"net/http"
	"strconv"
	"strings"
	
	"github.com/adyutaa/parsea/internal/service"
	"github.com/adyutaa/parsea/internal/validation"
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
	CVID     uint   `json:"cv_id" binding:"required"`
	ReportID uint   `json:"report_id" binding:"required"`
	JobTitle string `json:"job_title" binding:"required"`
}

// Evaluate creates a new evaluation job
func (h *EvaluationHandler) Evaluate(c *gin.Context) {
	var req EvaluateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON format: " + err.Error(),
		})
		return
	}

	// Validate CV ID (basic validation - ensure it's not zero)
	if req.CVID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cv_id must be a positive integer",
		})
		return
	}

	// Validate Report ID (basic validation - ensure it's not zero)
	if req.ReportID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "report_id must be a positive integer",
		})
		return
	}

	// Validate and sanitize job title
	req.JobTitle = strings.TrimSpace(req.JobTitle)
	if err := validation.ValidateJobTitle(req.JobTitle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Start evaluation
	jobID, err := h.service.StartEvaluation(
		strconv.FormatUint(uint64(req.CVID), 10), 
		strconv.FormatUint(uint64(req.ReportID), 10), 
		req.JobTitle,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start evaluation: " + err.Error(),
		})
		return
	}

	// Convert jobID string back to int for response
	jobIDInt, _ := strconv.Atoi(jobID)
	c.JSON(http.StatusOK, gin.H{
		"id":     jobIDInt,
		"status": "queued",
	})
}

// GetResult retrieves the result of an evaluation job
func (h *EvaluationHandler) GetResult(c *gin.Context) {
	jobID := c.Query("id")

	// Validate job ID format
	if err := validation.ValidateID(jobID, "id"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"hint": "Job ID must be a valid positive integer. Example: ?id=123",
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
