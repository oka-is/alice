package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/pkg/domain"
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

func TestStorage_UpdateCredentials(t *testing.T) {
	store, savepoint := MustNewStore(t)
	t.Cleanup(savepoint.Flush)

	t.Run("it update user credentials", func(t *testing.T) {
		ctx := context.Background()
		oldUser := mustCreateUser(t, store, &domain.User{})
		extUser := mustCreateUser(t, store, &domain.User{})
		newUser := domain.User{
			Identity:   domain.NewEmptyString("1"),
			Verifier:   domain.NewEmptyBytes([]byte{1}),
			SrpSalt:    domain.NewEmptyBytes([]byte{2}),
			PasswdSalt: domain.NewEmptyBytes([]byte{3}),
			PrivKeyEnc: domain.NewEmptyBytes([]byte{4}),
		}

		err := store.UpdateCredentials(ctx, oldUser.ID.String, oldUser.Identity.String, newUser)
		require.NoError(t, err)

		oldUser01, err01 := store.FindUser(ctx, oldUser.ID.String)
		extUser01, err02 := store.FindUser(ctx, extUser.ID.String)

		require.NoError(t, err01)
		require.NoError(t, err02)

		require.Equal(t, newUser.Identity.String, oldUser01.Identity.String, "Identity")
		require.Equal(t, newUser.Verifier.Bytea, oldUser01.Verifier.Bytea, "Verifier")
		require.Equal(t, newUser.SrpSalt.Bytea, oldUser01.SrpSalt.Bytea, "SrpSalt")
		require.Equal(t, newUser.PasswdSalt.Bytea, oldUser01.PasswdSalt.Bytea, "PasswdSalt")
		require.Equal(t, newUser.PrivKeyEnc.Bytea, oldUser01.PrivKeyEnc.Bytea, "PrivKeyEnc")

		require.Equal(t, extUser.Identity.String, extUser01.Identity.String, "Identity")
		require.Equal(t, extUser.Verifier.Bytea, extUser01.Verifier.Bytea, "Verifier")
		require.Equal(t, extUser.SrpSalt.Bytea, extUser01.SrpSalt.Bytea, "SrpSalt")
		require.Equal(t, extUser.PasswdSalt.Bytea, extUser01.PasswdSalt.Bytea, "PasswdSalt")
		require.Equal(t, extUser.PrivKeyEnc.Bytea, extUser01.PrivKeyEnc.Bytea, "PrivKeyEnc")
	})
}

func TestStorage_IssueUserOtp(t *testing.T) {
	store, savepoint := MustNewStore(t)
	t.Cleanup(savepoint.Flush)

	t.Run("it works", func(t *testing.T) {
		ctx := context.Background()
		user01 := mustCreateUser(t, store, &domain.User{})
		user02 := mustCreateUser(t, store, &domain.User{})
		secret := "foo"

		err := store.IssueUserOtp(ctx, user01.ID.String, secret)
		require.NoError(t, err)

		user11, err11 := store.FindUser(ctx, user01.ID.String)
		user22, err22 := store.FindUser(ctx, user02.ID.String)
		require.NoError(t, err11)
		require.NoError(t, err22)

		require.Equal(t, secret, string(user11.OtpCandidate.Bytea))
		require.Equal(t, user02.OtpCandidate.Bytea, user22.OtpCandidate.Bytea)
	})
}

func TestStorage_EnableUserOtp(t *testing.T) {
	store, savepoint := MustNewStore(t)
	t.Cleanup(savepoint.Flush)

	t.Run("it works", func(t *testing.T) {
		ctx := context.Background()
		user01 := mustCreateUser(t, store, &domain.User{OtpCandidate: domain.NewEmptyBytes([]byte{1})})
		user02 := mustCreateUser(t, store, &domain.User{})

		err := store.EnableUserOtp(ctx, user01.ID.String, user01.Identity.String, user01.OtpCandidate.Bytea)
		require.NoError(t, err)

		user11, err11 := store.FindUser(ctx, user01.ID.String)
		user22, err22 := store.FindUser(ctx, user02.ID.String)
		require.NoError(t, err11)
		require.NoError(t, err22)

		require.Equal(t, user01.OtpCandidate.Bytea, user11.OtpSecret.Bytea, "sets secret from candidate")
		require.Equal(t, []byte{}, user11.OtpCandidate.Bytea, "clear candidate")

		require.Equal(t, user02.OtpSecret.Bytea, user22.OtpSecret.Bytea)
	})
}

func TestStorage_DisableUserOtp(t *testing.T) {
	store, savepoint := MustNewStore(t)
	t.Cleanup(savepoint.Flush)

	t.Run("it works", func(t *testing.T) {
		ctx := context.Background()
		user01 := mustCreateUser(t, store, &domain.User{OtpSecret: domain.NewEmptyBytes([]byte{1})})
		user02 := mustCreateUser(t, store, &domain.User{OtpSecret: domain.NewEmptyBytes([]byte{2})})

		err := store.DisableUserOtp(ctx, user01.ID.String)
		require.NoError(t, err)

		user11, err11 := store.FindUser(ctx, user01.ID.String)
		user22, err22 := store.FindUser(ctx, user02.ID.String)
		require.NoError(t, err11)
		require.NoError(t, err22)

		require.Equal(t, []byte{}, user11.OtpSecret.Bytea)
		require.Equal(t, user02.OtpSecret.Bytea, user22.OtpSecret.Bytea)
	})
}

func mustBuildUser(t *testing.T, storage *Storage, input *domain.User) *domain.User {
	out := &domain.User{
		Ver:          domain.NewEmptyInt64(1),
		Identity:     domain.NewEmptyString(domain.NewUUID()),
		Verifier:     domain.NewEmptyBytes([]byte("Verifier")),
		SrpSalt:      domain.NewEmptyBytes([]byte("SrpSalt")),
		PasswdSalt:   domain.NewEmptyBytes([]byte("PasswdSalt")),
		PrivKeyEnc:   domain.NewEmptyBytes([]byte("PrivKeyEnc")),
		PubKey:       domain.NewEmptyBytes([]byte("PubKey")),
		OtpCandidate: domain.NewNullBytea(),
		OtpSecret:    domain.NewNullBytea(),
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

	if input.OtpCandidate.Valid {
		out.OtpCandidate = input.OtpCandidate
	}

	if input.OtpSecret.Valid {
		out.OtpSecret = input.OtpSecret
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
