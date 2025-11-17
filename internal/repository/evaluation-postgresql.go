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

// Create saves a new evaluation job using direct SQL (no prepared statements)
func (r *EvaluationRepository) Create(job *domain.EvaluationJob) error {
	// Get underlying sql.DB and use it directly
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	
	query := `INSERT INTO evaluation_jobs (id, cv_id, report_id, job_title, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err = sqlDB.Exec(query, job.ID, job.CVID, job.ReportID, job.JobTitle, job.Status, job.CreatedAt, job.UpdatedAt)
	return err
}

// GetByID retrieves an evaluation job by ID using RAW SQL
func (r *EvaluationRepository) GetByID(id string) (*domain.EvaluationJob, error) {
	var job domain.EvaluationJob
	sql := `SELECT id, cv_id, report_id, job_title, status, result, error_message, created_at, updated_at 
	        FROM evaluation_jobs WHERE id = $1 LIMIT 1`
	
	err := r.db.Raw(sql, id).Scan(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// GetDB returns the database connection for debugging
func (r *EvaluationRepository) GetDB() *gorm.DB {
	return r.db
}

// UpdateStatus updates the status of a job using direct SQL
func (r *EvaluationRepository) UpdateStatus(id string, status string) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	
	query := `UPDATE evaluation_jobs SET status = $1, updated_at = $2 WHERE id = $3`
	_, err = sqlDB.Exec(query, status, time.Now(), id)
	return err
}

// UpdateResult updates the job with final results using direct SQL
func (r *EvaluationRepository) UpdateResult(id string, result *domain.EvaluationResult) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	
	// Use proper JSON marshaling instead of manual formatting
	resultMap := map[string]interface{}{
		"cv_match_rate":    result.CVMatchRate,
		"cv_feedback":      result.CVFeedback,
		"project_score":    result.ProjectScore,
		"project_feedback": result.ProjectFeedback,
		"overall_summary":  result.OverallSummary,
	}
	
	// Marshal to JSON bytes
	resultJSON, err := json.Marshal(resultMap)
	if err != nil {
		return fmt.Errorf("failed to marshal result to JSON: %w", err)
	}
	
	query := `UPDATE evaluation_jobs SET result = $1, status = $2, updated_at = $3 WHERE id = $4`
	_, err = sqlDB.Exec(query, string(resultJSON), "completed", time.Now(), id)
	return err
}

// UpdateError updates the job with error message using direct SQL
func (r *EvaluationRepository) UpdateError(id string, errMsg string) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	
	query := `UPDATE evaluation_jobs SET status = $1, error_message = $2, updated_at = $3 WHERE id = $4`
	_, err = sqlDB.Exec(query, "failed", errMsg, time.Now(), id)
	return err
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
