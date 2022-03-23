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

func ErrorJSON(err error) interface{} {
	if err == nil {
		return ""
	}

	switch v := err.(type) {
	case validation.Errors:
		return toJSON(v)
	default:
		return ErrorJSON(errors.Unwrap(err))
	}
}

func toJSON(err validation.Errors) map[string]string {
	out := make(map[string]string)
	for field, exception := range err {
		out[field] = exception.Error()
	}
	return out
}
