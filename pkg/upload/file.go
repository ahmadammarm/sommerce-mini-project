package upload

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/lucsky/cuid"
)

const (
	MaxFileSize = 5 * 1024 * 1024 // 5 MB
)

// ValidateFile checks the size and extension of the uploaded image.
func ValidateFile(file *multipart.FileHeader) error {
	// Check size
	if file.Size > MaxFileSize {
		return errors.New("file size exceeds 5MB limit")
	}

	// Check extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return errors.New("only PNG and JPG/JPEG files are allowed")
	}

	return nil
}

// GenerateFilename creates a unique CUID-based filename to prevent collisions.
func GenerateFilename(originalName string) string {
	ext := strings.ToLower(filepath.Ext(originalName))
	return cuid.New() + ext
}
