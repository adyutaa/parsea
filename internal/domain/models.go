package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Document struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid"`
	Filename   string    `json:"filename" gorm:"not null"`
	FilePath   string    `json:"file_path" gorm:"not null"`
	DocType    string    `json:"doc_type" gorm:"not null"`
	FileSize   int64     `json:"file_size"`
	UploadedAt time.Time `json:"uploaded_at" gorm:"default:now()"`
}

func (Document) TableName() string {
	return "documents"
}

type JSON map[string]interface{}

func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}
	
	return json.Unmarshal(bytes, j)
}

type EvaluationJob struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid"`
	CVID         string    `json:"cv_id" gorm:"type:uuid;not null"`
	ReportID     string    `json:"report_id" gorm:"type:uuid;not null"`
	JobTitle     string    `json:"job_title" gorm:"not null"`
	Status       string    `json:"status" gorm:"default:'queued'"` // queued, processing, completed, failed
	Result       JSON      `json:"result,omitempty" gorm:"type:jsonb"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:now()"`
}

func (EvaluationJob) TableName() string {
	return "evaluation_jobs"
}

type EvaluationResult struct {
	CVMatchRate     float64 `json:"cv_match_rate"`
	CVFeedback      string  `json:"cv_feedback"`
	ProjectScore    float64 `json:"project_score"`
	ProjectFeedback string  `json:"project_feedback"`
	OverallSummary  string  `json:"overall_summary"`
}
