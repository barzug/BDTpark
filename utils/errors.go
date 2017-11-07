package utils

import "errors"

var (
	UniqueError = errors.New("uniqueError")
	NotFoundError = errors.New("notFoundError")
)
