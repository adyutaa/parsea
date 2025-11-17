package domain

import (
	"github.com/openai/openai-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// LLM Service Types

type OpenAIClient struct {
	Client *openai.Client
}

type CVEvaluationResult struct {
	MatchRate float64 `json:"match_rate"`
	Feedback  string  `json:"feedback"`
}

type ProjectEvaluationResult struct {
	Score    float64 `json:"score"`
	Feedback string  `json:"feedback"`
}

// Infrastructure Service Types

type RedisClient struct {
	Client *redis.Client
}

type DatabaseConnection struct {
	DB *gorm.DB
}

// File Upload Types

type FileUpload struct {
	Filename    string
	ContentType string
	Size        int64
	Data        []byte
}

type FileValidationResult struct {
	IsValid bool
	Errors  []string
}

// API Response Types

type UploadResponse struct {
	CVID     string `json:"cv_id"`
	ReportID string `json:"report_id"`
	Message  string `json:"message"`
}

type EvaluationResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type ResultResponse struct {
	ID        string      `json:"id"`
	Status    string      `json:"status"`
	Result    JSON        `json:"result,omitempty"`
	CreatedAt any `json:"created_at"`
	UpdatedAt any `json:"updated_at"`
}

type QueueStatusResponse struct {
	QueueLength int64  `json:"queue_length"`
	Status      string `json:"status"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
	Hint    string `json:"hint,omitempty"`
}

// Validation Types

type ValidationRule struct {
	Field     string
	Required  bool
	MaxLength int
	MinLength int
	Pattern   string
}

type ValidationResult struct {
	IsValid bool
	Errors  map[string]string
}