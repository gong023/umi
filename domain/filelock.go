package domain

// FileLock is an interface for file locking operations
type FileLock interface {
	// Lock acquires a lock on the file at the given path
	// If the lock cannot be acquired, it returns an error
	Lock(path string) error

	// Unlock releases the lock on the file at the given path
	// If the lock cannot be released, it returns an error
	Unlock(path string) error

	// WithLock executes the given function while holding a lock on the file at the given path
	// It acquires the lock before executing the function and releases it after the function returns
	// If the lock cannot be acquired, it returns an error
	WithLock(path string, fn func() error) error
}
