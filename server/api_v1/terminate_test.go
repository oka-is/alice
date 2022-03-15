package api_v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/server/policy"
)

func TestTerminate(t *testing.T) {
	t.Run("it errors when user is a readonly", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.store.EXPECT().RetrieveSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.Session{}, nil)
		s.store.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(domain.User{}, nil)
		s.userPolicy.EXPECT().Wrap(gomock.Any()).Return(s.userPolicy)
		s.userPolicy.EXPECT().CanWrite().Return(policy.ErrDenied)

		s.MustPOST(t, "/v1/terminate", &alice_v1.TerminateRequest{})

		require.Equal(t, 403, s.res.Code)
	})

	t.Run("it ok when user is ok", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.store.EXPECT().RetrieveSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.Session{}, nil)
		s.store.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(domain.User{}, nil)
		s.userPolicy.EXPECT().Wrap(gomock.Any()).Return(s.userPolicy)
		s.userPolicy.EXPECT().CanWrite().Return(nil)
		s.store.EXPECT().TerminateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		s.MustPOST(t, "/v1/terminate", &alice_v1.TerminateRequest{})

		require.Equal(t, 200, s.res.Code)
	})
}
