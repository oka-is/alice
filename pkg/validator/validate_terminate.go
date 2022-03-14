package validator

import (
	"bytes"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oka-is/alice/pkg/domain"
)

type ValidateTerminateOpts struct {
	User     domain.User
	Identity []byte
}

func (v *Validator) ValidateTerminate(opts ValidateTerminateOpts) error {
	return validation.Errors{
		"Identity": one(validation.Validate(bytes.Equal(opts.Identity, opts.User.Identity.Bytea), validation.Required)),
	}.Filter()
}

func (v NoOptValidator) ValidateTerminate(opts ValidateTerminateOpts) error {
	return nil
}
