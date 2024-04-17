package consts

import "errors"


const (
	ErrExternalCode = 400
	ErrInternalCode = 500
)

var (
	ErrInternal error = errors.New("internal error")
	ErrExternal error = errors.New("external error")
)