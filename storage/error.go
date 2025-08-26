package storage

import "fmt"

var ErrCouldNotAcquireLock = fmt.Errorf("could not acquire lock")
var ErrFailedToAcquireLock = fmt.Errorf("failed to acquire lock")
