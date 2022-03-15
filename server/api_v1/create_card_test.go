package api_v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/server/policy"
)

func TestCreateCard(t *testing.T) {
	t.Run("when user is a readonly", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})
		s.userPolicy.EXPECT().CanWrite().Return(policy.ErrDenied)

		s.MustPOST(t, "/v1/workspaces/:wid/cards/create", &alice_v1.UpsertCardRequest{})
		require.Equal(t, 403, s.res.Code)
	})

	t.Run("when unauthorized", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})
		s.userPolicy.EXPECT().CanWrite().Return(nil)
		s.workspacePolicy.EXPECT().CanManageWorkspace().Return(policy.ErrDenied)
		s.store.EXPECT().FindUserWorkspaceLink(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.UserWorkspace{}, nil)

		s.MustPOST(t, "/v1/workspaces/:wid/cards/create", &alice_v1.UpsertCardRequest{})
		require.Equal(t, 403, s.res.Code)
	})

	t.Run("when ok", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		card := domain.Card{TitleEnc: domain.NewEmptyBytes([]byte{1})}
		user := domain.User{ID: domain.NewEmptyString("foo")}

		s.LoginAs(t, domain.Session{}, user)
		s.userPolicy.EXPECT().CanWrite().Return(nil)
		s.workspacePolicy.EXPECT().CanManageWorkspace().Return(nil)
		s.store.EXPECT().FindUserWorkspaceLink(gomock.Any(), user.ID.String, ":wid").Return(domain.UserWorkspace{}, nil)
		s.store.EXPECT().CreateCardWithItems(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		s.MustPOST(t, "/v1/workspaces/:wid/cards/create", &alice_v1.UpsertCardRequest{
			Card: &alice_v1.Card{TitleEnc: card.TitleEnc.Bytea},
		})

		res := new(alice_v1.UpsertCardResponse)
		s.MustBindResponse(t, res)

		require.Equal(t, 200, s.res.Code)
		require.Equal(t, card.TitleEnc.Bytea, res.Card.TitleEnc)
	})
}
