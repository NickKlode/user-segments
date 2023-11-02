package service

import "fmt"

var (
	ErrUserAlreadyExists           = fmt.Errorf("user already exists")
	ErrCannotCreateUser            = fmt.Errorf("cannot crate user")
	ErrUserAlreadyInSegments       = fmt.Errorf("user already in segment")
	ErrCannotAddUser               = fmt.Errorf("cannot add user to segment")
	ErrUserNotFound                = fmt.Errorf("user not found")
	ErrUserNotInSegment            = fmt.Errorf("user not in segment")
	ErrCannotDeleteUserFromSegment = fmt.Errorf("cannot delete user from segment")
	ErrCannotGetUserSegmets        = fmt.Errorf("cannot get user segments")
	ErrNoSegments                  = fmt.Errorf("users has no segments")

	ErrSegmentAlreadyExists = fmt.Errorf("segment already exists")
	ErrCannotCreateSegment  = fmt.Errorf("cannot crate segment")
	ErrCannotDeleteSegment  = fmt.Errorf("cannot delete segment")
	ErrSegmentNotFound      = fmt.Errorf("segment not found")

	ErrOperationHistoryNotFound  = fmt.Errorf("operation history not found")
	ErrCannotGetOperationHistory = fmt.Errorf("cannot get operation history")
)
