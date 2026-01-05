package errs

import (
	"errors"
	"fmt"
)

var ErrValueIsInvalid = errors.New("value is invalid")

type ValueInvalidError struct {
	ParamName string
	Cause     error
}

func NewValueInvalidErrorWithReason(paramName string, reason error) *ValueInvalidError {
	return &ValueInvalidError{
		ParamName: paramName,
		Cause:     reason,
	}
}

func NewValueInvalidError(paramName string) *ValueInvalidError {
	return &ValueInvalidError{
		ParamName: paramName,
	}
}

func (e *ValueInvalidError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)",
			ErrValueIsRequired, e.ParamName, e.Cause,
		)
	}
	return fmt.Sprintf("%s: %s", ErrValueIsRequired, e.ParamName)
}

func (e *ValueInvalidError) Unwrap() error {
	return ErrValueIsInvalid
}
