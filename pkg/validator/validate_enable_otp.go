package validator

import (
	"bytes"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/wault-pw/alice/pkg/domain"
)

type ValidateEnableUserOtpOpts struct {
	User     domain.User
	Identity string
	Secret   []byte
}

func (v *Validator) ValidateEnableUserOtp(opts ValidateEnableUserOtpOpts) error {
	return validation.Errors{
		"Identity":     one(validation.Validate(opts.Identity == opts.User.Identity.String, validation.Required)),
		"OtpSecret":    one(validation.Validate(!opts.User.OtpSecret.Valid, validation.Required)),
		"OtpCandidate": one(validation.Validate(bytes.Equal(opts.User.OtpCandidate.Bytea, opts.Secret), validation.Required)),
	}.Filter()
}

func (v NoOptValidator) ValidateEnableUserOtp(opts ValidateEnableUserOtpOpts) error {
	return nil
}
