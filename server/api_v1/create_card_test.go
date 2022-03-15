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
	t.Run("it errors when user is a readonly", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})
		s.userPolicy.EXPECT().CanWrite().Return(policy.ErrDenied)

		s.MustPOST(t, "/v1/workspaces/:wid/cards/create", &alice_v1.UpsertCardRequest{})

		require.Equal(t, 403, s.res.Code)
	})

	t.Run("it ok", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		card := domain.Card{TitleEnc: domain.NewEmptyBytes([]byte{1})}

		s.LoginAs(t, domain.Session{}, domain.User{})
		s.userPolicy.EXPECT().CanWrite().Return(nil)
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
