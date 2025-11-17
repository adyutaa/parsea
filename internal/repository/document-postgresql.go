package repository

import (
	"github.com/adyutaa/parsea/internal/domain"

	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// Create saves a new document to the database
func (r *DocumentRepository) Create(doc *domain.Document) error {
	return r.db.Create(doc).Error
}

// GetByID retrieves a document by its ID
func (r *DocumentRepository) GetByID(id uint) (*domain.Document, error) {
	var doc domain.Document
	err := r.db.Where("id = ?", id).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// GetByType retrieves all documents of a specific type
func (r *DocumentRepository) GetByType(docType string) ([]domain.Document, error) {
	var docs []domain.Document
	err := r.db.Where("doc_type = ?", docType).Find(&docs).Error
	return docs, err
}
