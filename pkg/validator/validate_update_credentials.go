package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/wault-pw/alice/pkg/domain"
)

type ValidateUpdateCredentialsOpts struct {
	OldUser     domain.User
	NewUser     domain.User
	OldIdentity string
}

func (v *Validator) ValidateUpdateCredentials(opts ValidateUpdateCredentialsOpts) error {
	return validation.Errors{
		"OldIdentity": one(validation.Validate(opts.OldIdentity == opts.OldUser.Identity.String, validation.Required)),
		"Identity":    one(validation.Validate(opts.NewUser.Identity, validation.Required)),
		"Verifier":    one(validation.Validate(opts.NewUser.Verifier, validation.Required)),
		"SrpSalt":     one(validation.Validate(opts.NewUser.SrpSalt, validation.Required)),
		"PasswdSalt":  one(validation.Validate(opts.NewUser.PasswdSalt, validation.Required)),
		"PrivKeyEnc":  one(validation.Validate(opts.NewUser.PrivKeyEnc, validation.Required)),
	}.Filter()
}

func (v NoOptValidator) ValidateUpdateCredentials(opts ValidateUpdateCredentialsOpts) error {
	return nil
}
