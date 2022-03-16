package api_v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/server/policy"
)

func TestUpdateWorkspace(t *testing.T) {
	t.Run("when user is a readonly", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})
		s.userPolicy.EXPECT().CanWrite().Return(policy.ErrDenied)

		s.MustPOST(t, "/v1/workspaces/:wid/update", nil)

		require.Equal(t, 403, s.res.Code)
	})

	t.Run("when unauthorized", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})

		s.userPolicy.EXPECT().CanWrite().Return(nil)
		s.workspacePolicy.EXPECT().CanManageWorkspace().Return(policy.ErrDenied)
		s.store.EXPECT().FindUserWorkspaceLink(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.UserWorkspace{}, nil)

		s.MustPOST(t, "/v1/workspaces/:wid/update", nil)
		require.Equal(t, 403, s.res.Code)
	})

	t.Run("when ok", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		user := domain.User{ID: domain.NewEmptyString("foo")}
		uw := domain.UserWorkspace{ID: domain.NewEmptyString("UserWorkspaceID"), WorkspaceID: domain.NewEmptyString("bat")}
		userWithWorkspace := domain.UserWithWorkspace{TitleEnc: domain.NewEmptyBytes([]byte{1, 2, 3})}
		title := []byte{1}
		s.LoginAs(t, domain.Session{}, user)

		s.userPolicy.EXPECT().CanWrite().Return(nil)
		s.workspacePolicy.EXPECT().CanManageWorkspace().Return(nil)
		s.store.EXPECT().FindUserWorkspaceLink(gomock.Any(), user.ID.String, ":wid").Return(uw, nil)
		s.store.EXPECT().UpdateWorkspace(gomock.Any(), uw.WorkspaceID.String, title).Return(nil)
		s.store.EXPECT().FindUserWithWorkspace(gomock.Any(), uw.ID.String).Return(userWithWorkspace, nil)

		s.MustPOST(t, "/v1/workspaces/:wid/update", &alice_v1.UpdateWorkspaceRequest{TitleEnc: title})
		require.Equal(t, 200, s.res.Code)

		res := new(alice_v1.UpdateWorkspaceResponse)
		s.MustBindResponse(t, res)

		require.Equal(t, userWithWorkspace.TitleEnc.Bytea, res.Workspace.TitleEnc)
	})
}
