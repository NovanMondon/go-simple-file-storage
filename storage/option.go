package storage

import "time"

type Option func(*storageConfigs)

func WithLockPath(lockPath string) Option {
	return func(c *storageConfigs) {
		c.lockPath = lockPath
	}
}

func WithCheckInterval(interval time.Duration) Option {
	return func(c *storageConfigs) {
		c.checkInterval = interval
	}
}

func WithRetryMax(retryMax int) Option {
	return func(c *storageConfigs) {
		c.retryMax = retryMax
	}
}
