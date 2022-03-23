package validator

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	MaxLoginAttempts = 3
	LoginAttemptsDur = time.Minute
)

type ValidateCandidateSessionOpts struct {
	Attempts int
}

func (v *Validator) ValidateCandidateSession(opts ValidateCandidateSessionOpts) error {
	return validation.Errors{
		"Attempts": one(validation.Validate(opts.Attempts <= MaxLoginAttempts, validation.Required.Error("max login attempts, please try again later"))),
	}.Filter()
}

func (v NoOptValidator) ValidateCandidateSession(opts ValidateCandidateSessionOpts) error {
	return nil
}
