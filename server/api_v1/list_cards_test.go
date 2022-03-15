package api_v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/server/policy"
)

func TestListCards(t *testing.T) {
	t.Run("when unauthorized", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		s.LoginAs(t, domain.Session{}, domain.User{})

		s.workspacePolicy.EXPECT().CanSeeWorkspace().Return(policy.ErrDenied)
		s.store.EXPECT().FindUserWorkspaceLink(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.UserWorkspace{}, nil)

		s.MustPOST(t, "/v1/workspaces/:wid/cards", nil)
		require.Equal(t, 403, s.res.Code)
	})

	t.Run("when ok", func(t *testing.T) {
		s := MustSetup(t)
		defer s.ctrl.Finish()

		user := domain.User{ID: domain.NewEmptyString("foo")}
		uw := domain.UserWorkspace{WorkspaceID: domain.NewEmptyString("bar")}
		card := domain.Card{ID: domain.NewEmptyString("baz")}
		s.LoginAs(t, domain.Session{}, user)

		s.workspacePolicy.EXPECT().CanSeeWorkspace().Return(nil)
		s.store.EXPECT().FindUserWorkspaceLink(gomock.Any(), user.ID.String, ":wid").Return(uw, nil)
		s.store.EXPECT().ListCardsByWorkspace(gomock.Any(), uw.WorkspaceID.String).Return([]domain.Card{card}, nil)

		s.MustPOST(t, "/v1/workspaces/:wid/cards", nil)
		require.Equal(t, 200, s.res.Code)

		res := new(alice_v1.ListCardsResponse)
		s.MustBindResponse(t, res)

		require.Len(t, res.Items, 1)
		require.Equal(t, res.Items[0].Id, card.ID.String)
	})
}
