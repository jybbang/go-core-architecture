package core

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrNotFound            = errors.New("your requested Item is not found")
	ErrConflict            = errors.New("your Item already exist")
	ErrBadRequest          = errors.New("given Param is not valid")
	ErrForbiddenAcccess    = errors.New("your access is forbidden")
)
