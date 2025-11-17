package pdf

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// ExtractText extracts text content from a PDF file
func (p *Parser) ExtractText(filePath string) (string, error) {
	// Open PDF file
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	// Get total pages
	totalPages := r.NumPage()

	var textBuilder strings.Builder

	// Extract text from each page
	for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		text, err := p.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %w", pageIndex, err)
		}

		textBuilder.WriteString(text)
		textBuilder.WriteString("\n")
	}

	extractedText := textBuilder.String()

	// Validate extracted text
	if len(strings.TrimSpace(extractedText)) == 0 {
		return "", fmt.Errorf("no text content found in PDF")
	}

	return extractedText, nil
}

// ExtractTextFromBytes extracts text from PDF bytes
func (p *Parser) ExtractTextFromBytes(data []byte) (string, error) {
	r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("failed to read PDF: %w", err)
	}

	totalPages := r.NumPage()
	var textBuilder strings.Builder

	for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		text, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}

		textBuilder.WriteString(text)
		textBuilder.WriteString("\n")
	}

	return textBuilder.String(), nil
}

// CleanText removes extra whitespace and normalizes text
func (p *Parser) CleanText(text string) string {
	// Remove multiple spaces
	text = strings.Join(strings.Fields(text), " ")

	// Remove multiple newlines
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n")
}
