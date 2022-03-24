package validator

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// NewError used in testing purpose
func NewError(message string) error {
	return validation.Errors{"BASE": validation.NewError("INTERNAL", message)}.Filter()
}

func IsInvalid(err error) bool {
	return errors.As(err, &validation.Errors{})
}
