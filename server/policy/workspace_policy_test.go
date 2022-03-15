package policy

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/pkg/domain"
)

func TestWorkspacePolicy_CanManageWorkspace(t *testing.T) {
	type args struct {
		desc  string
		build func() (domain.User, domain.UserWorkspace)
		want  error
	}

	tests := []args{
		{
			desc: "when foreign owner",
			build: func() (domain.User, domain.UserWorkspace) {
				user, uw, _ := MustBuildWorkspaceManageCard(t)
				uw.OwnerID = domain.NewEmptyString("foo")
				return user, uw
			},
			want: ErrDenied,
		}, {
			build: func() (domain.User, domain.UserWorkspace) {
				user, uw, _ := MustBuildWorkspaceManageCard(t)
				return user, uw
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			policy := new(WorkspacePolicy)
			user, uw := tt.build()
			policy.Wrap(user, uw)
			require.Equal(t, tt.want, policy.CanManageWorkspace())
		})
	}
}

func TestWorkspacePolicy_CanManageCard(t *testing.T) {
	type args struct {
		desc  string
		build func() (domain.User, domain.UserWorkspace, domain.Card)
		want  error
	}

	tests := []args{
		{
			desc: "when foreign owner",
			build: func() (domain.User, domain.UserWorkspace, domain.Card) {
				user, uw, card := MustBuildWorkspaceManageCard(t)
				uw.OwnerID = domain.NewEmptyString("foo")
				return user, uw, card
			},
			want: ErrDenied,
		}, {
			desc: "when foreign card",
			build: func() (domain.User, domain.UserWorkspace, domain.Card) {
				user, uw, card := MustBuildWorkspaceManageCard(t)
				card.WorkspaceID = domain.NewEmptyString("foo")
				return user, uw, card
			},
			want: ErrDenied,
		}, {
			desc: "when ok",
			build: func() (domain.User, domain.UserWorkspace, domain.Card) {
				user, uw, card := MustBuildWorkspaceManageCard(t)
				return user, uw, card
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			policy := new(WorkspacePolicy)
			user, uw, card := tt.build()
			policy.Wrap(user, uw)
			require.Equal(t, tt.want, policy.CanManageCard(card))
		})
	}
}

func TestWorkspacePolicy_CanSeeWorkspace(t *testing.T) {
	type args struct {
		desc  string
		build func() (domain.User, domain.UserWorkspace)
		want  error
	}

	tests := []args{
		{
			desc: "when foreign user",
			build: func() (domain.User, domain.UserWorkspace) {
				user, uw, _ := MustBuildWorkspaceSeeCard(t)
				uw.UserID = domain.NewEmptyString("foo")
				return user, uw
			},
			want: ErrDenied,
		}, {
			build: func() (domain.User, domain.UserWorkspace) {
				user, uw, _ := MustBuildWorkspaceSeeCard(t)
				return user, uw
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			policy := new(WorkspacePolicy)
			user, uw := tt.build()
			policy.Wrap(user, uw)
			require.Equal(t, tt.want, policy.CanSeeWorkspace())
		})
	}
}

func TestWorkspacePolicy_CanSeeCard(t *testing.T) {
	type args struct {
		desc  string
		build func() (domain.User, domain.UserWorkspace, domain.Card)
		want  error
	}

	tests := []args{
		{
			desc: "when foreign user",
			build: func() (domain.User, domain.UserWorkspace, domain.Card) {
				user, uw, card := MustBuildWorkspaceSeeCard(t)
				uw.UserID = domain.NewEmptyString("foo")
				return user, uw, card
			},
			want: ErrDenied,
		}, {
			desc: "when foreign card",
			build: func() (domain.User, domain.UserWorkspace, domain.Card) {
				user, uw, card := MustBuildWorkspaceSeeCard(t)
				card.WorkspaceID = domain.NewEmptyString("foo")
				return user, uw, card
			},
			want: ErrDenied,
		}, {
			desc: "when ok",
			build: func() (domain.User, domain.UserWorkspace, domain.Card) {
				user, uw, card := MustBuildWorkspaceSeeCard(t)
				return user, uw, card
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			policy := new(WorkspacePolicy)
			user, uw, card := tt.build()
			policy.Wrap(user, uw)
			require.Equal(t, tt.want, policy.CanSeeCard(card))
		})
	}
}

// MustBuildWorkspaceManageCard builds resources with valid IDs
func MustBuildWorkspaceManageCard(t *testing.T) (domain.User, domain.UserWorkspace, domain.Card) {
	user := domain.User{ID: domain.NewEmptyString(domain.NewUUID())}
	uw := domain.UserWorkspace{
		OwnerID:     user.ID,
		WorkspaceID: domain.NewEmptyString(domain.NewUUID()),
	}
	card := domain.Card{WorkspaceID: uw.WorkspaceID}
	return user, uw, card
}

// MustBuildWorkspaceSeeCard builds resources with valid IDs
func MustBuildWorkspaceSeeCard(t *testing.T) (domain.User, domain.UserWorkspace, domain.Card) {
	user := domain.User{ID: domain.NewEmptyString(domain.NewUUID())}
	uw := domain.UserWorkspace{
		UserID:      user.ID,
		WorkspaceID: domain.NewEmptyString(domain.NewUUID()),
	}
	card := domain.Card{WorkspaceID: uw.WorkspaceID}
	return user, uw, card
}
