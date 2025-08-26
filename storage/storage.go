package storage

import (
	"os"
	"time"

	"github.com/gofrs/flock"
)

type Storage[T any] struct {
	c storageConfigs

	marshal   func(T) ([]byte, error)
	unmarshal func([]byte) (T, error)

	lock *flock.Flock
}

type storageConfigs struct {
	filePath      string
	lockPath      string
	checkInterval time.Duration
	retryMax      int
}

func New[T any](
	filePath string,
	marshal func(T) ([]byte, error),
	unmarshal func([]byte) (T, error),
	opts ...Option,
) Storage[T] {
	storage := Storage[T]{
		c: storageConfigs{
			filePath:      filePath,
			lockPath:      filePath + ".lock",
			checkInterval: 100 * time.Millisecond,
			retryMax:      -1,
		},

		marshal:   marshal,
		unmarshal: unmarshal,
	}

	for _, opt := range opts {
		opt(&storage.c)
	}

	storage.lock = flock.New(storage.c.lockPath)

	return storage
}

func (s Storage[T]) TryLoad() (T, error) {
	var none T
	locked, err := s.lock.TryLock()
	if err != nil {
		return none, err
	}
	if !locked {
		return none, ErrCouldNotAcquireLock
	}
	defer s.lock.Unlock()

	data, err := os.ReadFile(s.c.filePath)
	if err != nil {
		return none, err
	}

	result, err := s.unmarshal(data)
	if err != nil {
		return none, err
	}

	return result, nil
}

func (s Storage[T]) TrySave(content T) error {
	locked, err := s.lock.TryLock()
	if err != nil {
		return err
	}
	if !locked {
		return ErrCouldNotAcquireLock
	}
	defer s.lock.Unlock()

	data, err := s.marshal(content)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.c.filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage[T]) Load() (T, error) {
	retryCount := 0
	for {
		content, err := s.TryLoad()
		if err == ErrCouldNotAcquireLock {
			time.Sleep(s.c.checkInterval)
			retryCount++
			if s.c.retryMax >= 0 && retryCount > s.c.retryMax {
				break
			}
			continue
		}
		return content, err
	}

	var none T
	return none, ErrFailedToAcquireLock
}

func (s Storage[T]) Save(content T) error {
	retry := 0
	for {
		err := s.TrySave(content)
		if err == ErrCouldNotAcquireLock {
			time.Sleep(s.c.checkInterval)
			retry++
			if s.c.retryMax >= 0 && retry > s.c.retryMax {
				break
			}
			continue
		}
		return err
	}

	return ErrFailedToAcquireLock
}

func (s Storage[T]) TryOpen() (OpenedStorage[T], error) {
	locked, err := s.lock.TryLock()
	if err != nil {
		return OpenedStorage[T]{}, err
	}
	if !locked {
		return OpenedStorage[T]{}, ErrCouldNotAcquireLock
	}

	return OpenedStorage[T]{
		filePath:  s.c.filePath,
		marshal:   s.marshal,
		unmarshal: s.unmarshal,
	}, nil
}

func (s Storage[T]) Open() (OpenedStorage[T], error) {
	retryCount := 0
	for {
		file, err := s.TryOpen()
		if err == ErrCouldNotAcquireLock {
			time.Sleep(s.c.checkInterval)
			retryCount++
			if s.c.retryMax >= 0 && retryCount > s.c.retryMax {
				break
			}
			continue
		}
		return file, err
	}

	return OpenedStorage[T]{}, ErrFailedToAcquireLock
}

func (s Storage[T]) Close() error {
	return s.lock.Unlock()
}
