package storage

import (
	"os"
	"time"

	"github.com/gofrs/flock"
)

type StringStorage struct {
	filePath      string
	lockPath      string
	checkInterval time.Duration

	lock *flock.Flock
}

func NewStringStorage(filePath string, opts ...Option) StringStorage {
	// use default values
	storage := StringStorage{
		filePath:      filePath,
		lockPath:      filePath + ".lock",
		checkInterval: 100 * time.Millisecond,
	}

	// apply options
	for _, opt := range opts {
		opt(&storage)
	}

	storage.lock = flock.New(storage.lockPath)

	return storage
}

func (s StringStorage) Save(data string) error {
	for {
		err := s.TrySave(data)
		if err == ErrCouldNotAcquireLock {
			time.Sleep(s.checkInterval)
			continue
		}
		return err
	}
}

func (s StringStorage) Load() (string, error) {
	for {
		data, err := s.TryLoad()
		if err == ErrCouldNotAcquireLock {
			time.Sleep(s.checkInterval)
			continue
		}
		return data, err
	}
}

func (s StringStorage) TrySave(data string) error {
	locked, err := s.lock.TryLock()
	if err != nil {
		return err
	}
	if !locked {
		return ErrCouldNotAcquireLock
	}

	if err := os.WriteFile(s.filePath, []byte(data), 0644); err != nil {
		return err
	}

	if err := s.lock.Unlock(); err != nil {
		return err
	}

	return nil
}

func (s StringStorage) TryLoad() (string, error) {
	locked, err := s.lock.TryLock()
	if err != nil {
		return "", err
	}
	if !locked {
		return "", ErrCouldNotAcquireLock
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return "", err
	}

	if err := s.lock.Unlock(); err != nil {
		return "", err
	}

	return string(data), nil
}
