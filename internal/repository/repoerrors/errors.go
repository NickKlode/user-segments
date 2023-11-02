package repoerrors

import "errors"

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrSegmentNotFound          = errors.New("segment not found")
	ErrAlreadyExists            = errors.New("already exists")
	ErrUserNotInSegment         = errors.New("user not in segment")
	ErrOperationHistoryNotFound = errors.New("operation history not found")
)
