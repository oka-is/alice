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
		require.Equal(t, user.Verifier.Bytea, user1.Verifier.Bytea, "id decrypts encrypted column")
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

func TestStorage_TerminateUser(t *testing.T) {
	store, savepoint := MustNewStore(t)
	t.Cleanup(savepoint.Flush)

	t.Run("it deletes user and all its resources", func(t *testing.T) {
		ctx := context.Background()
		user := mustCreateUser(t, store, &domain.User{})

		workspace := mustCreateWorkspace(t, store, &domain.Workspace{})
		mustCreateUserWorkspace(t, store, &domain.UserWorkspace{
			OwnerID:     user.ID,
			WorkspaceID: workspace.ID,
		})

		err := store.TerminateUser(ctx, user.Identity.String, user.ID.String)
		require.NoError(t, err)

		_, err01 := store.FindUser(ctx, user.ID.String)
		_, err02 := store.FindWorkspace(ctx, workspace.ID.String)

		require.ErrorIs(t, ErrNotFound, err01)
		require.ErrorIs(t, ErrNotFound, err02)
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
