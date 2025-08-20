package util

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInternal      = errors.New("internal error")
)
