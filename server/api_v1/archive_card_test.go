package api_v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/server/policy"
)

func TestArchiveCard(t *testing.T) {
	t.Run("when user is a readonly", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})
		s.userPolicy.EXPECT().CanWrite().Return(policy.ErrDenied)

		s.MustPOST(t, "/v1/workspaces/foo/cards/bar/archive", nil)
		require.Equal(t, 403, s.res.Code)
	})

	t.Run("when unauthorized", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})
		s.userPolicy.EXPECT().CanWrite().Return(nil)
		s.workspacePolicy.EXPECT().CanManageCard(gomock.Any()).Return(policy.ErrDenied)
		s.store.EXPECT().FindUserWorkspaceLink(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.UserWorkspace{}, nil)
		s.store.EXPECT().FindCard(gomock.Any(), gomock.Any()).Return(domain.Card{}, nil)

		s.MustPOST(t, "/v1/workspaces/foo/cards/bar/archive", nil)
		require.Equal(t, 403, s.res.Code)
	})

	t.Run("when ok", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		user := domain.User{ID: domain.NewEmptyString("foo")}
		s.LoginAs(t, domain.Session{}, user)
		s.userPolicy.EXPECT().CanWrite().Return(nil)
		s.workspacePolicy.EXPECT().CanManageCard(gomock.Any()).Return(nil)
		s.store.EXPECT().FindUserWorkspaceLink(gomock.Any(), user.ID.String, ":wid").Return(domain.UserWorkspace{}, nil)
		s.store.EXPECT().FindCard(gomock.Any(), ":cid").Return(domain.Card{}, nil)
		s.store.EXPECT().ArchiveCard(gomock.Any(), gomock.Any()).Return(true, nil)

		s.MustPOST(t, "/v1/workspaces/:wid/cards/:cid/archive", nil)
		res := new(alice_v1.ArchiveCardResponse)
		s.MustBindResponse(t, res)

		require.Equal(t, 200, s.res.Code)
		require.Equal(t, true, res.Archived)
	})
}
