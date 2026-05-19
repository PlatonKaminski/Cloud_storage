package domain

import (
	"errors"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrFileNotFound       = errors.New("file not found")
	ErrFileAccessDenied   = errors.New("access denied")
	ErrFileTooLarge       = errors.New("file too large")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("wrong login/password")
	ErrAlreadyExists      = errors.New("already exists")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrFolderAccessDenied = errors.New("folder access denied")
	ErrFolderNotFound     = errors.New("folder not found")
)
