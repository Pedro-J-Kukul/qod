// Filename: internal/data/errors.go
// Description: Custom error definitions for data operations
package data

import (
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

var ErrEditConflict = errors.New("edit conflict")

var ErrDuplicateEmail = errors.New("duplicate email")

var ErrInvalidCredentials = errors.New("invalid credentials")

var ErrNoRecords = errors.New("no matching records found")
