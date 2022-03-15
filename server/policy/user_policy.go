package policy

import "github.com/wault-pw/alice/pkg/domain"

type UserPolicy struct {
	user domain.User
}

func (u *UserPolicy) Wrap(user domain.User) IUserPolicy {
	u.user = user
	return u
}

// CanWrite cant modify ane resource even self profile
func (u *UserPolicy) CanWrite() error {
	if u.user.Readonly.Bool {
		return ErrDenied
	}
	return nil
}
