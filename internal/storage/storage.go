package storage

import "errors"

var (
	ErrUrlExists     = errors.New("url exists")
	ErrAliasNotFount = errors.New("alias not found")
	ErrUrlNotFound   = errors.New("url not found")
)
