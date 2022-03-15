package policy

import (
	"errors"
	"github.com/wault-pw/alice/pkg/domain"
)

var ErrDenied = errors.New("ERR_DENIED")

//go:generate mockgen -destination ../policy_mock/policy_mock.go -source types.go -package policy_mock -mock_names IUserPolicy=MockUserPolicy
type IUserPolicy interface {
	Wrap(user domain.User) IUserPolicy
	CanWrite() error
}
