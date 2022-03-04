package storage

import (
	"context"
	"testing"

	"github.com/oka-is/alice/pkg/domain"
)

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
