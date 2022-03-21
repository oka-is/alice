package api_v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/server/policy"
)

func TestOtpEnable(t *testing.T) {
	t.Run("when user is a readonly", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})
		s.userPolicy.EXPECT().CanWrite().Return(policy.ErrDenied)

		s.MustPOST(t, "/v1/otp/enable", nil)

		require.Equal(t, 403, s.res.Code)
	})

	t.Run("when ok", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		passcode := s.opts.OtpStub
		user := domain.User{
			Identity:     domain.NewEmptyString("foo"),
			ID:           domain.NewEmptyString("bar"),
			OtpCandidate: domain.NewEmptyBytes([]byte{1}),
		}
		s.LoginAs(t, domain.Session{}, user)

		s.userPolicy.EXPECT().CanWrite().Return(nil)
		s.store.EXPECT().EnableUserOtp(gomock.Any(), user.ID.String, user.Identity.String, user.OtpCandidate.Bytea)

		s.MustPOST(t, "/v1/otp/enable", &alice_v1.OtpEnableRequest{
			Identity: user.Identity.String,
			Passcode: passcode,
		})

		require.Equal(t, 200, s.res.Code)
	})
}
