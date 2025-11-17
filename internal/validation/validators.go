package validation

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func IsValidID(s string) bool {
	if s == "" {
		return false
	}
	_, err := strconv.ParseUint(s, 10, 32)
	return err == nil
}

func ValidateID(s string, fieldName string) error {
	if s == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	if !IsValidID(s) {
		return fmt.Errorf("%s must be a valid positive integer", fieldName)
	}
	return nil
}

func ValidateJobTitle(jobTitle string) error {
	if jobTitle == "" {
		return fmt.Errorf("job_title is required")
	}

	jobTitle = strings.TrimSpace(jobTitle)
	if len(jobTitle) == 0 {
		return fmt.Errorf("job_title cannot be empty")
	}

	if len(jobTitle) > 100 {
		return fmt.Errorf("job_title cannot exceed 100 characters")
	}

	validPattern := regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)
	if !validPattern.MatchString(jobTitle) {
		return fmt.Errorf("job_title contains invalid characters")
	}

	return nil
}

func ValidateFilename(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename is required")
	}

	if len(filename) > 255 {
		return fmt.Errorf("filename too long (max 255 characters)")
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".pdf" {
		return fmt.Errorf("only PDF files are allowed")
	}

	return nil
}

func ValidateFileSize(size int64) error {
	const maxSize = 10 * 1024 * 1024

	if size <= 0 {
		return fmt.Errorf("file is empty")
	}

	if size > maxSize {
		return fmt.Errorf("file size exceeds 10MB limit")
	}

	return nil
}

func ValidateMimeType(mimeType string) error {
	allowedTypes := []string{
		"application/pdf",
	}

	for _, allowed := range allowedTypes {
		if mimeType == allowed {
			return nil
		}
	}

	return fmt.Errorf("invalid file type: %s", mimeType)
}
