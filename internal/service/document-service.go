package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"github.com/adyutaa/parsea/internal/domain"
	"github.com/adyutaa/parsea/internal/repository"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type DocumentService struct {
	repo       *repository.DocumentRepository
	uploadPath string
}

func NewDocumentService(repo *repository.DocumentRepository, uploadPath string) *DocumentService {
	return &DocumentService{
		repo:       repo,
		uploadPath: uploadPath,
	}
}

// SaveDocument saves an uploaded file and stores its metadata
func (s *DocumentService) SaveDocument(file *multipart.FileHeader, docType string) (string, error) {
	// Validate file type
	ext := filepath.Ext(file.Filename)
	if ext != ".pdf" {
		return "", fmt.Errorf("only PDF files are allowed")
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		return "", fmt.Errorf("file size exceeds 10MB limit")
	}

	// Generate unique ID
	id := uuid.New().String()

	// Create file path
	filename := fmt.Sprintf("%s%s", id, ext)
	filePath := filepath.Join(s.uploadPath, filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Save metadata to database
	doc := &domain.Document{
		ID:         id,
		Filename:   file.Filename,
		FilePath:   filePath,
		DocType:    docType,
		FileSize:   file.Size,
		UploadedAt: time.Now(),
	}

	if err := s.repo.Create(doc); err != nil {
		// Clean up file if database save fails
		os.Remove(filePath)
		return "", fmt.Errorf("failed to save document metadata: %w", err)
	}

	return id, nil
}

// GetDocument retrieves document metadata by ID
func (s *DocumentService) GetDocument(id string) (*domain.Document, error) {
	return s.repo.GetByID(id)
}
