package api_v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/server/policy"
)

func TestCloneCard(t *testing.T) {
	t.Run("when user is a readonly", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})
		s.userPolicy.EXPECT().CanWrite().Return(policy.ErrDenied)

		s.MustPOST(t, "/v1/workspaces/foo/cards/bar/clone", &alice_v1.CloneCardRequest{})
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

		s.MustPOST(t, "/v1/workspaces/:wid/cards/:cid/clone", &alice_v1.CloneCardRequest{})
		require.Equal(t, 403, s.res.Code)
	})

	t.Run("when ok", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		card := domain.Card{ID: domain.NewEmptyString("foo"), TitleEnc: domain.NewEmptyBytes([]byte{1})}
		user := domain.User{ID: domain.NewEmptyString("bar")}

		s.LoginAs(t, domain.Session{}, user)
		s.userPolicy.EXPECT().CanWrite().Return(nil)
		s.workspacePolicy.EXPECT().CanManageCard(gomock.Any()).Return(nil)
		s.store.EXPECT().FindUserWorkspaceLink(gomock.Any(), user.ID.String, ":wid").Return(domain.UserWorkspace{}, nil)
		s.store.EXPECT().FindCard(gomock.Any(), ":cid").Return(card, nil)
		s.store.EXPECT().CloneCard(gomock.Any(), card.ID.String, gomock.Any()).Return(card, nil)

		s.MustPOST(t, "/v1/workspaces/:wid/cards/:cid/clone", &alice_v1.CloneCardRequest{})
		res := new(alice_v1.CloneCardResponse)
		s.MustBindResponse(t, res)

		require.Equal(t, 200, s.res.Code)
		require.Equal(t, card.TitleEnc.Bytea, res.Card.TitleEnc)
	})
}
