package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// ValidateUploadedFile validates the extension, size, and MIME type of an uploaded file.
func ValidateUploadedFile(fileHeader *multipart.FileHeader, allowedExts map[string]bool, maxBytes int64, allowedMimeTypes map[string]bool) error {
	if fileHeader == nil {
		return errors.New("file is nil")
	}

	// 1. Validate extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExts[ext] {
		return errors.New("unsupported file extension")
	}

	// 2. Validate size
	if fileHeader.Size > maxBytes {
		return fmt.Errorf("file size exceeds the limit of %d MB", maxBytes/(1<<20))
	}

	// 3. Validate MIME type
	file, err := fileHeader.Open()
	if err != nil {
		return errors.New("cannot open file")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return errors.New("cannot read file")
	}

	mimeType := http.DetectContentType(buffer)
	if !allowedMimeTypes[mimeType] {
		return errors.New("invalid MIME type")
	}

	return nil
}

// SaveUploadedFileWithRandomName generates a random UUID filename with the original extension,
// ensures the destination directory exists, saves the file, and returns the relative image URL path.
func SaveUploadedFileWithRandomName(fileHeader *multipart.FileHeader, saveDir string) (string, error) {
	if fileHeader == nil {
		return "", errors.New("file is nil")
	}

	// 1. Generate unique filename (UUID + original extension)
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 2. Ensure save directory exists locally under the working directory
	fullSaveDir := filepath.Join(GetWorkingDir(), saveDir)
	if err := os.MkdirAll(fullSaveDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// 3. Open source file
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()

	// 4. Create destination file
	dstPath := filepath.Join(fullSaveDir, uniqueName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// 5. Copy content
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return the relative path (e.g., "uploads/profile-images/uuid.png")
	return filepath.Join(saveDir, uniqueName), nil
}
