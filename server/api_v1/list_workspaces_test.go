package api_v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
)

func TestListWorkspaces(t *testing.T) {
	t.Run("when ok", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		user := domain.User{ID: domain.NewEmptyString("foo")}
		uw := domain.UserWithWorkspace{WorkspaceID: domain.NewEmptyString("bar")}

		s.LoginAs(t, domain.Session{}, user)
		s.store.EXPECT().ListUserWithWorkspaces(gomock.Any(), user.ID.String).Return([]domain.UserWithWorkspace{uw}, nil)

		s.MustPOST(t, "/v1/workspaces/list", nil)
		require.Equal(t, 200, s.res.Code)

		res := new(alice_v1.ListWorkspacesResponse)
		s.MustBindResponse(t, res)

		require.Len(t, res.Items, 1)
		require.Equal(t, res.Items[0].WorkspaceId, uw.WorkspaceID.String)
	})
}
