package pdf

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ledongthuc/pdf"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ExtractText(filePath string) (string, error) {
	// Try method 1: ledongthuc/pdf library
	text, err := p.extractWithLedong(filePath)
	if err == nil && len(strings.TrimSpace(text)) > 0 {
		return text, nil
	}

	// Try method 2: pdftotext (if available)
	text, err = p.extractWithPdfToText(filePath)
	if err == nil && len(strings.TrimSpace(text)) > 0 {
		return text, nil
	}

	// If both methods fail
	return "", fmt.Errorf("no text content found in PDF")
}

func (p *Parser) extractWithLedong(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	totalPages := r.NumPage()
	var textBuilder strings.Builder

	for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
		page := r.Page(pageIndex)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue // Skip failed pages
		}

		textBuilder.WriteString(text)
		textBuilder.WriteString("\n")
	}

	return textBuilder.String(), nil
}

func (p *Parser) extractWithPdfToText(filePath string) (string, error) {
	// Check if pdftotext is available
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return "", fmt.Errorf("pdftotext not available: %w", err)
	}

	// Run pdftotext command
	cmd := exec.Command("pdftotext", filePath, "-")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("pdftotext failed: %w", err)
	}

	return string(output), nil
}

func (p *Parser) ExtractTextFromBytes(data []byte) (string, error) {
	r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("failed to read PDF: %w", err)
	}

	totalPages := r.NumPage()
	var textBuilder strings.Builder

	for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
		page := r.Page(pageIndex)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}

		textBuilder.WriteString(text)
		textBuilder.WriteString("\n")
	}

	return textBuilder.String(), nil
}

func (p *Parser) CleanText(text string) string {

	text = strings.Join(strings.Fields(text), " ")

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
