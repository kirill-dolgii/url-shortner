package storage

import "errors"

var (
	ErrUrlExists = errors.New("url exists")
	ErrNotFound  = errors.New("not found")
)
