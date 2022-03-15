package api_v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/server/policy"
)

func TestCreateWorkspace(t *testing.T) {
	t.Run("it errors when user is a readonly", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})

		s.userPolicy.EXPECT().CanWrite().Return(policy.ErrDenied)
		s.MustPOST(t, "/v1/workspaces/create", &alice_v1.CreateWorkspaceRequest{})

		require.Equal(t, 403, s.res.Code)
	})

	t.Run("it ok", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})

		workspace := domain.UserWithWorkspace{
			TitleEnc: domain.NewEmptyBytes([]byte{1}),
		}

		s.userPolicy.EXPECT().CanWrite().Return(nil)
		s.store.EXPECT().CreateWorkspace(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.store.EXPECT().FindUserWithWorkspace(gomock.Any(), gomock.Any()).Return(workspace, nil)

		s.MustPOST(t, "/v1/workspaces/create", &alice_v1.CreateWorkspaceRequest{})

		res := new(alice_v1.CreateWorkspaceResponse)
		s.MustBindResponse(t, res)

		require.Equal(t, 200, s.res.Code)
		require.Equal(t, workspace.TitleEnc.Bytea, res.GetWorkspace().GetTitleEnc())
	})
}
