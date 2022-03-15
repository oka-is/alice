package api_v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/server/policy"
)

func TestDeleteCard(t *testing.T) {
	t.Run("it errors when user is a readonly", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})
		s.userPolicy.EXPECT().CanWrite().Return(policy.ErrDenied)

		s.MustPOST(t, "/v1/workspaces/:wid/cards/:cid/delete", nil)

		require.Equal(t, 403, s.res.Code)
	})

	t.Run("it ok", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})
		s.userPolicy.EXPECT().CanWrite().Return(nil)
		s.store.EXPECT().DeleteCard(gomock.Any(), gomock.Any()).Return(nil)

		s.MustPOST(t, "/v1/workspaces/:wid/cards/:cid/delete", nil)

		require.Equal(t, 200, s.res.Code)
	})
}
