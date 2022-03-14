package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/wault-pw/alice/pkg/domain"
)

type ValidateTerminateOpts struct {
	User     domain.User
	Identity string
}

func (v *Validator) ValidateTerminate(opts ValidateTerminateOpts) error {
	return validation.Errors{
		"Identity": one(validation.Validate(opts.Identity == opts.User.Identity.String, validation.Required)),
	}.Filter()
}

func (v NoOptValidator) ValidateTerminate(opts ValidateTerminateOpts) error {
	return nil
}
