package usecase

import (
	"path/filepath"
)

// MockFilepathHandler is a helper struct to handle filepath operations in tests
type MockFilepathHandler struct {
	TempDir string
}

// NewMockFilepathHandler creates a new MockFilepathHandler
func NewMockFilepathHandler(tempDir string) *MockFilepathHandler {
	return &MockFilepathHandler{
		TempDir: tempDir,
	}
}

// Join is a wrapper around filepath.Join that redirects memo paths to the temp directory
func (h *MockFilepathHandler) Join(elem ...string) string {
	if len(elem) > 0 && elem[0] == "memo" {
		return filepath.Join(append([]string{h.TempDir}, elem...)...)
	}
	return filepath.Join(elem...)
}
