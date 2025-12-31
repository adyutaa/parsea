package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adyutaa/parsea/internal/domain"
	"gorm.io/gorm"
)

type EvaluationRepository struct {
	db *gorm.DB
}

func NewEvaluationRepository(db *gorm.DB) *EvaluationRepository {
	return &EvaluationRepository{db: db}
}

// Create saves a new evaluation job
func (r *EvaluationRepository) Create(job *domain.EvaluationJob) error {
	return r.db.Create(job).Error
}

// GetByID retrieves an evaluation job by ID
func (r *EvaluationRepository) GetByID(id uint) (*domain.EvaluationJob, error) {
	var job domain.EvaluationJob
	err := r.db.Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// GetDB returns the database connection for debugging
func (r *EvaluationRepository) GetDB() *gorm.DB {
	return r.db
}

// UpdateStatus updates the status of a job
func (r *EvaluationRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&domain.EvaluationJob{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// UpdateResult updates the job with final results
func (r *EvaluationRepository) UpdateResult(id string, result *domain.EvaluationResult) error {
	resultMap := map[string]interface{}{
		"cv_match_rate":    result.CVMatchRate,
		"cv_feedback":      result.CVFeedback,
		"project_score":    result.ProjectScore,
		"project_feedback": result.ProjectFeedback,
		"overall_summary":  result.OverallSummary,
	}
	
	resultJSON, err := json.Marshal(resultMap)
	if err != nil {
		return fmt.Errorf("failed to marshal result to JSON: %w", err)
	}
	
	return r.db.Model(&domain.EvaluationJob{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"result":     string(resultJSON),
			"status":     "completed",
			"updated_at": time.Now(),
		}).Error
}

// UpdateError updates the job with error message
func (r *EvaluationRepository) UpdateError(id string, errMsg string) error {
	return r.db.Model(&domain.EvaluationJob{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        "failed",
			"error_message": errMsg,
			"updated_at":    time.Now(),
		}).Error
}

// GetPendingJobs retrieves all jobs with status "queued"
func (r *EvaluationRepository) GetPendingJobs(limit int) ([]domain.EvaluationJob, error) {
	var jobs []domain.EvaluationJob
	err := r.db.Where("status = ?", "queued").
		Order("created_at ASC").
		Limit(limit).
		Find(&jobs).Error
	return jobs, err
}
