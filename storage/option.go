package storage

import "time"

type Option func(*StringStorage)

func WithLockPath(lockPath string) Option {
	return func(s *StringStorage) {
		s.lockPath = lockPath
	}
}

func WithCheckInterval(interval time.Duration) Option {
	return func(s *StringStorage) {
		s.checkInterval = interval
	}
}
