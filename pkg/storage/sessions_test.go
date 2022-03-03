package storage

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/oka-is/alice/lib/jwt_mock"
	"github.com/stretchr/testify/require"
)

func TestStorage_IssueSession(t *testing.T) {
	t.Run("it works", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ots := jwt_mock.NewMockOts(ctrl)
		defer ctrl.Finish()

		store, savepoint := MustNewStore(t)
		t.Cleanup(savepoint.Flush)

		ots.EXPECT().Marshall().Return("foo", nil)
		ots.EXPECT().SetJti(gomock.Any()).Return(ots)
		ots.EXPECT().SetExp(gomock.Any()).Return(ots)

		session, token, err := store.IssueSession(context.Background(), ots)
		require.NoError(t, err)
		require.Equal(t, "foo", token)
		require.NotEmpty(t, session.Jti.String)
		require.Equal(t, true, session.TimeTo.Time.After(session.TimeFrom.Time))
		require.Empty(t, session.UserID.String)
		require.Empty(t, session.CandidateID.String)
	})
}
