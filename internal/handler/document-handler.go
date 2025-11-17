package handler

import (
	"mime/multipart"
	"net/http"
	
	"github.com/adyutaa/parsea/internal/service"
	"github.com/adyutaa/parsea/internal/validation"
	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	service *service.DocumentService
}

func NewDocumentHandler(service *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

// Upload handles file uploads for CV and Project Report
func (h *DocumentHandler) Upload(c *gin.Context) {
	// Get CV file
	cvFile, err := c.FormFile("cv")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "CV file is required",
		})
		return
	}

	// Validate CV file
	if err := h.validateFile(cvFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "CV validation failed: " + err.Error(),
		})
		return
	}

	// Get Project Report file
	reportFile, err := c.FormFile("project_report")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Project report file is required",
		})
		return
	}

	// Validate Project Report file
	if err := h.validateFile(reportFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Project report validation failed: " + err.Error(),
		})
		return
	}

	// Save CV
	cvID, err := h.service.SaveDocument(cvFile, "cv")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save CV: " + err.Error(),
		})
		return
	}

	// Save Project Report
	reportID, err := h.service.SaveDocument(reportFile, "project_report")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save project report: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cv_id":     cvID,
		"report_id": reportID,
		"message":   "Files uploaded successfully",
	})
}

// validateFile performs comprehensive file validation
func (h *DocumentHandler) validateFile(file *multipart.FileHeader) error {
	// Validate filename
	if err := validation.ValidateFilename(file.Filename); err != nil {
		return err
	}

	// Validate file size
	if err := validation.ValidateFileSize(file.Size); err != nil {
		return err
	}

	// Validate MIME type
	if err := validation.ValidateMimeType(file.Header.Get("Content-Type")); err != nil {
		return err
	}

	return nil
}
