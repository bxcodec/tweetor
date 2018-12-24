package domain

import "errors"

var (
	ErrContextNil = errors.New("Context is Nil")
	ErrNotFound   = errors.New("Item doesn't exist")
)
