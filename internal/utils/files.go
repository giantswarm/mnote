package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetOutputPath generates the output path for a given input file and suffix
func GetOutputPath(inputPath, suffix string) string {
	dir := filepath.Dir(inputPath)
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	if suffix != "" {
		return filepath.Join(dir, fmt.Sprintf("%s_%s.md", base, suffix))
	}
	return filepath.Join(dir, base+".md")
}

// FileExists checks if a file exists and is not a directory
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// EnsureDirectory ensures that the directory for the given file path exists
func EnsureDirectory(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}

// WriteFile writes data to a file, creating the directory if needed
func WriteFile(path string, data []byte) error {
	if err := EnsureDirectory(path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// ReadFile reads the entire file into memory
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
