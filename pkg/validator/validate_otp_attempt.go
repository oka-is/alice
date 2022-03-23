package validator

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	MaxOtpAttempts = 3
	OtpAttemptsDur = time.Minute
)

type ValidateOtpAttemptOpts struct {
	Attempts int
}

func (v *Validator) ValidateOtpAttempt(opts ValidateOtpAttemptOpts) error {
	return validation.Errors{
		"Attempts": one(validation.Validate(opts.Attempts <= MaxOtpAttempts, validation.Required.Error("max otp attempts, please try again later"))),
	}.Filter()
}

func (v NoOptValidator) ValidateOtpAttempt(opts ValidateOtpAttemptOpts) error {
	return nil
}
