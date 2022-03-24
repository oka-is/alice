package mapper_v1

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/wault-pw/alice/desc/alice_v1"
)

func MapValidationError(err error) *alice_v1.ValidationError {
	if err == nil {
		return &alice_v1.ValidationError{}
	}

	switch e := err.(type) {
	case validation.Errors:
		return mapValidationError(e)
	default:
		return MapValidationError(errors.Unwrap(err))
	}
}

func mapValidationError(err validation.Errors) *alice_v1.ValidationError {
	out := &alice_v1.ValidationError{Items: make([]*alice_v1.ValidationError_Item, 0)}
	for field, exception := range err {
		out.Items = append(out.Items, &alice_v1.ValidationError_Item{
			Field:       field,
			Description: exception.Error(),
		})
	}

	return out
}
