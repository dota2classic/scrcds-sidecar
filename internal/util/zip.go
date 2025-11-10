package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// CompressFile compresses a single file into a ZIP archive with maximum compression.
func CompressFile(filePath string, zipPath string) error {
	// Read the source file into memory
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Create an in-memory buffer to hold the ZIP
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Set up a compressed entry
	w, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:   filepath.Base(filePath),
		Method: zip.Deflate, // enable compression
	})
	if err != nil {
		return fmt.Errorf("failed to create zip entry: %w", err)
	}

	// Write the file content into the ZIP entry
	if _, err := w.Write(fileData); err != nil {
		return fmt.Errorf("failed to write zip content: %w", err)
	}

	// Finalize the ZIP before using it
	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("failed to finalize zip: %w", err)
	}

	// Write the buffer to disk
	if err := os.WriteFile(zipPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write zip file to disk: %w", err)
	}

	return nil
}
