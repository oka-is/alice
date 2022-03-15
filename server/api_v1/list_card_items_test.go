package api_v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/server/policy"
)

func TestListCardItems(t *testing.T) {
	t.Run("when unauthorized", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})

		s.workspacePolicy.EXPECT().CanSeeCard(gomock.Any()).Return(policy.ErrDenied)
		s.store.EXPECT().FindUserWorkspaceLink(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.UserWorkspace{}, nil)
		s.store.EXPECT().FindCard(gomock.Any(), gomock.Any()).Return(domain.Card{}, nil)

		s.MustPOST(t, "/v1/workspaces/:wid/cards/:cid/items", nil)
		require.Equal(t, 403, s.res.Code)
	})

	t.Run("when ok", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		user := domain.User{ID: domain.NewEmptyString("foo")}
		card := domain.Card{ID: domain.NewEmptyString("bar")}
		item := domain.CardItem{ID: domain.NewEmptyString("baz")}
		uw := domain.UserWorkspace{WorkspaceID: domain.NewEmptyString("qwe")}
		s.LoginAs(t, domain.Session{}, user)

		s.workspacePolicy.EXPECT().CanSeeCard(gomock.Any()).Return(nil)
		s.store.EXPECT().FindUserWorkspaceLink(gomock.Any(), user.ID.String, ":wid").Return(uw, nil)
		s.store.EXPECT().FindCard(gomock.Any(), ":cid").Return(card, nil)
		s.store.EXPECT().ListCardItems(gomock.Any(), card.ID.String).Return([]domain.CardItem{item}, nil)

		s.MustPOST(t, "/v1/workspaces/:wid/cards/:cid/items", nil)
		require.Equal(t, 200, s.res.Code)

		res := new(alice_v1.ListCardItemsResponse)
		s.MustBindResponse(t, res)

		require.Len(t, res.Items, 1)
		require.Equal(t, res.Items[0].Id, item.ID.String)
	})
}
