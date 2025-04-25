package domain

// FileSystem is an interface for file system operations
type FileSystem interface {
	// ReadFile reads the content of a file at the given path
	ReadFile(path string) ([]byte, error)

	// WriteFile writes data to a file at the given path
	WriteFile(path string, data []byte, perm int) error

	// FileExists checks if a file exists at the given path
	FileExists(path string) (bool, error)

	// RemoveFile removes a file at the given path
	RemoveFile(path string) error

	// JoinPath joins path elements into a single path
	JoinPath(elem ...string) string
}
