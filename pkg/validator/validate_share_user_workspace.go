package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/wault-pw/alice/pkg/domain"
)

type ValidateShareUserWorkspaceOpts struct {
	NotShared     bool
	UserExists    bool
	UserWorkspace domain.UserWorkspace
}

func (v *Validator) ValidateShareUserWorkspace(opts ValidateShareUserWorkspaceOpts) error {
	return validation.Errors{
		"AedKeyEnc": one(validation.Validate(opts.UserWorkspace.AedKeyEnc, validation.Required)),
		"User": one(
			validation.Validate(opts.UserWorkspace.UserID != opts.UserWorkspace.OwnerID, validation.Required.Error("can't share with self")),
			validation.Validate(opts.NotShared, validation.Required.Error("is currently shared")),
			validation.Validate(opts.UserExists, validation.Required.Error("not exists")),
		),
	}.Filter()
}

func (v NoOptValidator) ValidateShareUserWorkspace(opts ValidateShareUserWorkspaceOpts) error {
	return nil
}
