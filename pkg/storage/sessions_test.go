package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/oka-is/alice/lib/jwt_mock"
	"github.com/oka-is/alice/pkg/domain"
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

func TestStorage_RetrieveSession(t *testing.T) {
	store, savepoint := MustNewStore(t)
	t.Cleanup(savepoint.Flush)

	type args struct {
		desc    string
		before  func(ots *jwt_mock.MockOts) *domain.Session
		wantJIT string
		wantErr bool
	}

	tests := []args{
		{
			desc: "when JIT valid",
			before: func(ots *jwt_mock.MockOts) *domain.Session {
				session := mustCreateSession(t, store, &domain.Session{})
				ots.EXPECT().Unmarshall(gomock.Any()).Return(session.Jti.String, nil)
				return session
			},
			wantErr: false,
		}, {
			desc: "when JIT invalid",
			before: func(ots *jwt_mock.MockOts) *domain.Session {
				ots.EXPECT().Unmarshall(gomock.Any()).Return("", errors.New("foo"))
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ots := jwt_mock.NewMockOts(ctrl)
			defer ctrl.Finish()

			want := tt.before(ots)
			session, err := store.RetrieveSession(context.Background(), ots, "")

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, want.Jti.String, session.Jti.String)
		})
	}
}

func mustBuildSession(t *testing.T, storage *Storage, input *domain.Session) *domain.Session {
	out := &domain.Session{
		Jti:      domain.NewEmptyString(domain.NewUUID()),
		TimeFrom: domain.NewEmptyTime(time.Now()),
		TimeTo:   domain.NewEmptyTime(time.Now().Add(24 * time.Hour)),
	}

	if input.Jti.Valid {
		out.Jti = input.Jti
	}

	if input.TimeFrom.Valid {
		out.TimeFrom = input.TimeFrom
	}

	if input.TimeTo.Valid {
		out.TimeTo = input.TimeTo
	}

	return out
}

func mustCreateSession(t *testing.T, storage *Storage, input *domain.Session) *domain.Session {
	ctx := context.Background()
	output := mustBuildSession(t, storage, input)
	if err := storage.insertSession(ctx, storage.db, output); err != nil {
		t.Errorf("cant create factory session: %s", err)
		t.FailNow()
	}
	return output
}
