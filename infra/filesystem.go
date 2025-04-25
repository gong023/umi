package infra

import (
	"os"
	"path/filepath"

	"github.com/gong023/umi/domain"
)

// FileSystem is an implementation of the domain.FileSystem interface
type FileSystem struct {
	logger domain.Logger
}

// NewFileSystem creates a new FileSystem instance
func NewFileSystem(logger domain.Logger) *FileSystem {
	return &FileSystem{
		logger: logger,
	}
}

// ReadFile reads the content of a file at the given path
func (fs *FileSystem) ReadFile(path string) ([]byte, error) {
	fs.logger.Info("Reading file: %s", path)
	content, err := os.ReadFile(path)
	if err != nil {
		fs.logger.Error("Failed to read file: %v", err)
		return nil, err
	}
	return content, nil
}

// WriteFile writes data to a file at the given path
func (fs *FileSystem) WriteFile(path string, data []byte, perm int) error {
	fs.logger.Info("Writing file: %s", path)
	
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fs.logger.Error("Failed to create directory: %v", err)
		return err
	}
	
	if err := os.WriteFile(path, data, os.FileMode(perm)); err != nil {
		fs.logger.Error("Failed to write file: %v", err)
		return err
	}
	return nil
}

// FileExists checks if a file exists at the given path
func (fs *FileSystem) FileExists(path string) (bool, error) {
	fs.logger.Info("Checking if file exists: %s", path)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	fs.logger.Error("Failed to check if file exists: %v", err)
	return false, err
}

// RemoveFile removes a file at the given path
func (fs *FileSystem) RemoveFile(path string) error {
	fs.logger.Info("Removing file: %s", path)
	if err := os.Remove(path); err != nil {
		fs.logger.Error("Failed to remove file: %v", err)
		return err
	}
	return nil
}

// JoinPath joins path elements into a single path
func (fs *FileSystem) JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}
