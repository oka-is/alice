package api_v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/server/policy"
)

func TestDeleteWorkspace(t *testing.T) {
	t.Run("when user is a readonly", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})
		s.userPolicy.EXPECT().CanWrite().Return(policy.ErrDenied)

		s.MustPOST(t, "/v1/workspaces/:wid/delete", nil)
		require.Equal(t, 403, s.res.Code)
	})

	t.Run("when unauthorized", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})

		s.userPolicy.EXPECT().CanWrite().Return(nil)
		s.workspacePolicy.EXPECT().CanManageWorkspace().Return(policy.ErrDenied)
		s.store.EXPECT().FindUserWorkspaceLink(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.UserWorkspace{}, nil)

		s.MustPOST(t, "/v1/workspaces/:wid/delete", nil)
		require.Equal(t, 403, s.res.Code)
	})

	t.Run("when ok", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		user := domain.User{ID: domain.NewEmptyString("foo")}
		uw := domain.UserWorkspace{WorkspaceID: domain.NewEmptyString("bat")}
		s.LoginAs(t, domain.Session{}, user)

		s.userPolicy.EXPECT().CanWrite().Return(nil)
		s.workspacePolicy.EXPECT().CanManageWorkspace().Return(nil)
		s.store.EXPECT().FindUserWorkspaceLink(gomock.Any(), user.ID.String, ":wid").Return(uw, nil)
		s.store.EXPECT().DeleteWorkspace(gomock.Any(), uw.WorkspaceID.String).Return(nil)

		s.MustPOST(t, "/v1/workspaces/:wid/delete", nil)
		require.Equal(t, 200, s.res.Code)
	})
}
