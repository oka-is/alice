package storage

import (
	"context"
	"testing"

	"github.com/oka-is/alice/pkg/domain"
	"github.com/stretchr/testify/require"
)

func TestStorage_FindUser(t *testing.T) {
	store, savepoint := MustNewStore(t)
	t.Cleanup(savepoint.Flush)

	t.Run("it works", func(t *testing.T) {
		ctx := context.Background()
		user := mustCreateUser(t, store, &domain.User{})

		user1, err1 := store.FindUser(ctx, user.ID.String)
		_, err2 := store.FindUser(ctx, domain.NewUUID())
		require.NoError(t, err1)
		require.ErrorIs(t, err2, ErrNotFound)

		require.Equal(t, user.ID.String, user1.ID.String)
	})
}

func TestStorage_FindUserIdentity(t *testing.T) {
	store, savepoint := MustNewStore(t)
	t.Cleanup(savepoint.Flush)

	t.Run("it works", func(t *testing.T) {
		ctx := context.Background()
		user := mustCreateUser(t, store, &domain.User{})

		user1, err1 := store.FindUserIdentity(ctx, user.Identity.String)
		_, err2 := store.FindUserIdentity(ctx, "foo")
		require.NoError(t, err1)
		require.ErrorIs(t, err2, ErrNotFound)

		require.Equal(t, user.ID.String, user1.ID.String)
	})
}

func mustBuildUser(t *testing.T, storage *Storage, input *domain.User) *domain.User {
	out := &domain.User{
		Ver:        domain.NewEmptyInt64(1),
		Identity:   domain.NewEmptyString(domain.NewUUID()),
		Verifier:   domain.NewEmptyBytes([]byte("Verifier")),
		SrpSalt:    domain.NewEmptyBytes([]byte("SrpSalt")),
		PasswdSalt: domain.NewEmptyBytes([]byte("PasswdSalt")),
		PrivKeyEnc: domain.NewEmptyBytes([]byte("PrivKeyEnc")),
		PubKey:     domain.NewEmptyBytes([]byte("PubKey")),
	}

	if input.Ver.Valid {
		out.Ver = input.Ver
	}

	if input.Identity.Valid {
		out.Identity = input.Identity
	}

	if input.Verifier.Valid {
		out.Verifier = input.Verifier
	}

	if input.SrpSalt.Valid {
		out.SrpSalt = input.SrpSalt
	}

	if input.PasswdSalt.Valid {
		out.PasswdSalt = input.PasswdSalt
	}

	if input.PrivKeyEnc.Valid {
		out.PrivKeyEnc = input.PrivKeyEnc
	}

	if input.PubKey.Valid {
		out.PubKey = input.PubKey
	}

	return out
}

func mustCreateUser(t *testing.T, storage *Storage, input *domain.User) *domain.User {
	ctx := context.Background()
	output := mustBuildUser(t, storage, input)
	if err := storage.insertUser(ctx, storage.db, output); err != nil {
		t.Errorf("cant create factory user: %s", err)
		t.FailNow()
	}
	return output
}
