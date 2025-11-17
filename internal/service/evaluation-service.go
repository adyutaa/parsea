package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adyutaa/parsea/internal/domain"
	"github.com/adyutaa/parsea/internal/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type EvaluationService struct {
	repo    *repository.EvaluationRepository
	docRepo *repository.DocumentRepository
	redis   *redis.Client
}

func NewEvaluationService(repo *repository.EvaluationRepository, docRepo *repository.DocumentRepository, redis *redis.Client) *EvaluationService {
	return &EvaluationService{
		repo:    repo,
		docRepo: docRepo,
		redis:   redis,
	}
}

func (s *EvaluationService) StartEvaluation(cvID, reportID, jobTitle string) (string, error) {
	// Convert string IDs to uint
	cvIDUint, err := strconv.ParseUint(cvID, 10, 32)
	if err != nil {
		return "", fmt.Errorf("invalid CV ID format: %w", err)
	}

	reportIDUint, err := strconv.ParseUint(reportID, 10, 32)
	if err != nil {
		return "", fmt.Errorf("invalid Report ID format: %w", err)
	}

	if _, err := s.docRepo.GetByID(uint(cvIDUint)); err != nil {
		return "", fmt.Errorf("CV document not found: %w", err)
	}

	if _, err := s.docRepo.GetByID(uint(reportIDUint)); err != nil {
		return "", fmt.Errorf("report document not found: %w", err)
	}

	job := &domain.EvaluationJob{
		CVID:      uint(cvIDUint),
		ReportID:  uint(reportIDUint),
		JobTitle:  jobTitle,
		Status:    "queued",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(job); err != nil {
		return "", fmt.Errorf("failed to create job: %w", err)
	}

	// PUSH REDIS queue
	ctx := context.Background()
	jobIDStr := strconv.FormatUint(uint64(job.ID), 10)
	if err := s.redis.LPush(ctx, "evaluation_queue", jobIDStr).Err(); err != nil {
		s.repo.UpdateError(jobIDStr, "failed to queue job")
		return "", fmt.Errorf("failed to queue job: %w", err)
	}

	return jobIDStr, nil
}

func (s *EvaluationService) GetJobStatus(id string) (*domain.EvaluationJob, error) {
	// Convert string ID to uint
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid job ID format: %w", err)
	}

	job, err := s.repo.GetByID(uint(idUint))
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}
	return job, nil
}

func (s *EvaluationService) GetDB() *gorm.DB {
	return s.repo.GetDB()
}

func (s *EvaluationService) GetQueueLength() (int64, error) {
	ctx := context.Background()
	length, err := s.redis.LLen(ctx, "evaluation_queue").Result()
	if err != nil {
		return 0, err
	}
	return length, nil
}
