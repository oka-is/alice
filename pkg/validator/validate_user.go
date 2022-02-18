package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oka-is/alice/pkg/domain"
)

func (v *Validator) ValidateUser(u domain.User) error {
	return validation.Errors{
		"Ver":        one(validation.Validate(u.Ver, validation.Required)),
		"Identity":   one(validation.Validate(u.Identity, validation.Required)),
		"Verifier":   one(validation.Validate(u.Verifier, validation.Required)),
		"SrpSalt":    one(validation.Validate(u.SrpSalt, validation.Required)),
		"PasswdSalt": one(validation.Validate(u.PasswdSalt, validation.Required)),
		"PrivKeyEnc": one(validation.Validate(u.PrivKeyEnc, validation.Required)),
		"PubKey":     one(validation.Validate(u.PubKey, validation.Required)),
	}.Filter()
}

func (v NoOptValidator) ValidateUser(u domain.User) error {
	return nil
}
