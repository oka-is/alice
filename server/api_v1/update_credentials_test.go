package api_v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/server/policy"
)

func TestUpdateCredentials(t *testing.T) {
	t.Run("when user is a readonly", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})
		s.userPolicy.EXPECT().CanWrite().Return(policy.ErrDenied)

		s.MustPOST(t, "/v1/credentials/update", nil)
		require.Equal(t, 403, s.res.Code)
	})

	t.Run("when ok", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		oldIdentity := "oldIdentity"
		user := domain.User{ID: domain.NewEmptyString("ID")}
		session := domain.Session{Jti: domain.NewEmptyString("JTI")}
		s.LoginAs(t, session, user)

		s.userPolicy.EXPECT().CanWrite().Return(nil)
		s.store.EXPECT().UpdateCredentials(gomock.Any(), user.ID.String, oldIdentity, gomock.Any()).Return(nil)
		s.store.EXPECT().DeleteUserSessionExcept(gomock.Any(), user.ID.String, session.Jti.String).Return(nil)
		s.store.EXPECT().FindUser(gomock.Any(), user.ID.String).Return(user, nil)

		s.MustPOST(t, "/v1/credentials/update", &alice_v1.UpdateCredentialsRequest{
			OldIdentity: oldIdentity,
		})

		res := new(alice_v1.UpdateCredentialsResponse)
		s.MustBindResponse(t, res)

		require.Equal(t, 200, s.res.Code)
		require.Equal(t, user.ID.String, res.User.Id)
	})
}
