package bus

import "errors"

var (
	ErrHandlerNotFound = errors.New("handler not found")
	ErrNonPointer      = errors.New("must be a non-nil pointer")
)
