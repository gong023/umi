package infra

import (
	"fmt"
	"sync"

	"github.com/gong023/umi/domain"
)

// FileLock is an implementation of the domain.FileLock interface
type FileLock struct {
	locks  map[string]*sync.Mutex
	mutex  sync.Mutex
	logger domain.Logger
}

// NewFileLock creates a new FileLock instance
func NewFileLock(logger domain.Logger) *FileLock {
	return &FileLock{
		locks:  make(map[string]*sync.Mutex),
		mutex:  sync.Mutex{},
		logger: logger,
	}
}

// Lock acquires a lock on the file at the given path
func (fl *FileLock) Lock(path string) error {
	fl.logger.Info("Acquiring lock for file: %s", path)
	
	// Get or create a mutex for the file
	fl.mutex.Lock()
	fileMutex, exists := fl.locks[path]
	if !exists {
		fileMutex = &sync.Mutex{}
		fl.locks[path] = fileMutex
	}
	fl.mutex.Unlock()
	
	// Acquire the lock
	fileMutex.Lock()
	fl.logger.Info("Lock acquired for file: %s", path)
	
	return nil
}

// Unlock releases the lock on the file at the given path
func (fl *FileLock) Unlock(path string) error {
	fl.logger.Info("Releasing lock for file: %s", path)
	
	// Get the mutex for the file
	fl.mutex.Lock()
	fileMutex, exists := fl.locks[path]
	fl.mutex.Unlock()
	
	if !exists {
		err := fmt.Errorf("no lock exists for file: %s", path)
		fl.logger.Error("Failed to release lock: %v", err)
		return err
	}
	
	// Release the lock
	fileMutex.Unlock()
	fl.logger.Info("Lock released for file: %s", path)
	
	return nil
}

// WithLock executes the given function while holding a lock on the file at the given path
func (fl *FileLock) WithLock(path string, fn func() error) error {
	// Acquire the lock
	if err := fl.Lock(path); err != nil {
		return err
	}
	
	// Ensure the lock is released even if the function panics
	defer func() {
		if err := fl.Unlock(path); err != nil {
			fl.logger.Error("Failed to release lock in defer: %v", err)
		}
	}()
	
	// Execute the function
	return fn()
}
