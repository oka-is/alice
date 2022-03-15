package policy

//go:generate mockgen -destination ../policy_mock/policy_mock.go -source types.go -package policy_mock -mock_names IUserPolicy=MockUserPolicy,IWorkspacePolicy=MockWorkspacePolicy

import (
	"errors"

	"github.com/wault-pw/alice/pkg/domain"
)

var ErrDenied = errors.New("ERR_DENIED")

type IUserPolicy interface {
	Wrap(user domain.User) IUserPolicy
	CanWrite() error
}

type IWorkspacePolicy interface {
	Wrap(user domain.User, uw domain.UserWorkspace) IWorkspacePolicy
	CanManageWorkspace() error
	CanManageCard(card domain.Card) error
	CanSeeWorkspace() error
	CanSeeCard(card domain.Card) error
}
