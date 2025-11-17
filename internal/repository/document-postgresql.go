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

// Create saves a new document to the database using simple protocol (no prepared statements)
func (r *DocumentRepository) Create(doc *domain.Document) error {
	// Get underlying sql.DB and use it directly to avoid any GORM magic
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	
	// Use direct SQL execution with simple protocol
	query := `INSERT INTO documents (id, filename, file_path, doc_type, file_size, uploaded_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	
	_, err = sqlDB.Exec(query, doc.ID, doc.Filename, doc.FilePath, doc.DocType, doc.FileSize, doc.UploadedAt)
	return err
}

// GetByID retrieves a document by its ID using RAW SQL
func (r *DocumentRepository) GetByID(id string) (*domain.Document, error) {
	var doc domain.Document
	sql := `SELECT id, filename, file_path, doc_type, file_size, uploaded_at 
	        FROM documents WHERE id = $1 LIMIT 1`
	
	err := r.db.Raw(sql, id).Scan(&doc).Error
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
