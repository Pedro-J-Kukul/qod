// Filename: internal/data/errors.go
// Description: Custom error definitions for data operations
package data

import (
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)
