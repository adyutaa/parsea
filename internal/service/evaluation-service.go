package service

import (
	"context"
	"fmt"
	"time"

	"github.com/adyutaa/parsea/internal/domain"
	"github.com/adyutaa/parsea/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9" // ← Update ini
	"gorm.io/gorm"
)

type EvaluationService struct {
	repo  *repository.EvaluationRepository
	redis *redis.Client
}

func NewEvaluationService(repo *repository.EvaluationRepository, redis *redis.Client) *EvaluationService {
	return &EvaluationService{
		repo:  repo,
		redis: redis,
	}
}

// StartEvaluation creates a new evaluation job and queues it
func (s *EvaluationService) StartEvaluation(cvID, reportID, jobTitle string) (string, error) {
	// Create job
	job := &domain.EvaluationJob{
		ID:        uuid.New().String(),
		CVID:      cvID,
		ReportID:  reportID,
		JobTitle:  jobTitle,
		Status:    "queued",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	if err := s.repo.Create(job); err != nil {
		return "", fmt.Errorf("failed to create job: %w", err)
	}

	// Push to Redis queue
	ctx := context.Background()
	if err := s.redis.LPush(ctx, "evaluation_queue", job.ID).Err(); err != nil {
		// If Redis fails, mark job as failed in DB
		s.repo.UpdateError(job.ID, "failed to queue job")
		return "", fmt.Errorf("failed to queue job: %w", err)
	}

	return job.ID, nil
}

// GetJobStatus retrieves the status of an evaluation job
func (s *EvaluationService) GetJobStatus(id string) (*domain.EvaluationJob, error) {
	job, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}
	return job, nil
}

// GetDB returns the database connection for debugging
func (s *EvaluationService) GetDB() *gorm.DB {
	return s.repo.GetDB()
}

// GetQueueLength returns the number of jobs in the queue
func (s *EvaluationService) GetQueueLength() (int64, error) {
	ctx := context.Background()
	length, err := s.redis.LLen(ctx, "evaluation_queue").Result()
	if err != nil {
		return 0, err
	}
	return length, nil
}
