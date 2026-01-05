package errs

import (
	"errors"
	"fmt"
)

var ErrValueIsRequired = errors.New("value is required")

type ValueRequiredError struct {
	ParamName string
	Cause     error
}

func NewValueRequiredErrorWithReason(paramName string, reason error) *ValueRequiredError {
	return &ValueRequiredError{
		ParamName: paramName,
		Cause:     reason,
	}
}

func NewValueRequiredError(paramName string) *ValueRequiredError {
	return &ValueRequiredError{
		ParamName: paramName,
	}
}

func (e *ValueRequiredError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)",
			ErrValueIsRequired, e.ParamName, e.Cause,
		)
	}
	return fmt.Sprintf("%s: %s", ErrValueIsRequired, e.ParamName)
}

func (e *ValueRequiredError) Unwrap() error {
	return ErrValueIsRequired
}
