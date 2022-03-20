package storage

import (
	"context"
	"fmt"

	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/pkg/validator"
)

func (s *Storage) CreateUser(ctx context.Context, user *domain.User, uw *domain.UserWorkspace,
	workspace *domain.Workspace, cardsWithItems []domain.CardWithItems) error {
	return s.Tx(ctx, nil, func(c context.Context, tx IConn) error {
		return s.createUser(c, tx, user, uw, workspace, cardsWithItems)
	})
}

func (s *Storage) FindUserIdentity(ctx context.Context, identity string) (user domain.User, err error) {
	query := s.selectUserColumns().From("users").Where("identity = ?", identity).Limit(1)
	err = s.Get(ctx, s.db, &user, query)
	return
}

func (s *Storage) FindUser(ctx context.Context, ID string) (user domain.User, err error) {
	return s.findUser(ctx, s.db, ID)
}

func (s *Storage) TerminateUser(ctx context.Context, identity string, userID string) error {
	return s.Tx(ctx, nil, func(c context.Context, tx IConn) error {
		return s.terminateUser(c, tx, identity, userID)
	})
}

func (s *Storage) UpdateCredentials(ctx context.Context, ID string, oldIdentity string, newUser domain.User) error {
	return s.Tx(ctx, nil, func(c context.Context, tx IConn) error {
		return s.updateCredentials(c, tx, ID, oldIdentity, newUser)
	})
}

func (s *Storage) IssueUserOtp(ctx context.Context, ID string, secret string) error {
	query := Builder().
		Update("users").
		Set("otp_candidate", Expr("encrypt(?, ?, 'aes-cbc/pad:pkcs')", secret, s.sseKey)).
		Where("id = ?", ID)

	return s.Exec1(ctx, s.db, query)
}

func (s *Storage) EnableUserOtp(ctx context.Context, ID string, identity string, secret []byte) error {
	return s.Tx(ctx, nil, func(c context.Context, tx IConn) error {
		return s.enableUserOtp(c, tx, ID, identity, secret)
	})
}

func (s *Storage) DisableUserOtp(ctx context.Context, ID string) error {
	query := Builder().Update("users").Set("otp_secret", Expr("NULL")).Where("id = ?", ID)
	return s.Exec1(ctx, s.db, query)
}

func (s *Storage) enableUserOtp(ctx context.Context, db IConn, ID string, identity string, secret []byte) error {
	user, err := s.findUser(ctx, db, ID)
	if err != nil {
		return fmt.Errorf("failed to find a user: %w", err)
	}

	err = s.validator.ValidateEnableUserOtp(validator.ValidateEnableUserOtpOpts{
		User:     user,
		Identity: identity,
		Secret:   secret,
	})
	if err != nil {
		return fmt.Errorf("failed to validate enable otp: %w", err)
	}

	query := Builder().Update("users").
		Set("otp_secret", Expr("otp_candidate")).
		Set("otp_candidate", Expr("NULL")).
		Where("id = ?", ID)

	return s.Exec1(ctx, db, query)

}

func (s *Storage) terminateUser(ctx context.Context, db IConn, identity string, userID string) error {
	user, err := s.findUser(ctx, db, userID)
	if err != nil {
		return fmt.Errorf("failed to find a user: %w", err)
	}

	err = s.validator.ValidateTerminate(validator.ValidateTerminateOpts{
		User:     user,
		Identity: identity,
	})
	if err != nil {
		return fmt.Errorf("failed to validate termination: %w", err)
	}

	err = s.deleteOwnerWorkspaces(ctx, db, userID)
	if err != nil {
		return fmt.Errorf("failed to delete workspaces: %w", err)
	}

	query := Builder().Delete("users").Where("id = ?", userID)
	return s.Exec1(ctx, db, query)
}

func (s *Storage) findUser(ctx context.Context, db IConn, ID string) (user domain.User, err error) {
	query := s.selectUserColumns().From("users").Where("id = ?", ID).Limit(1)
	err = s.Get(ctx, db, &user, query)
	return
}

func (s *Storage) createUser(ctx context.Context, db IConn, user *domain.User, uw *domain.UserWorkspace, workspace *domain.Workspace, cardsWithItems []domain.CardWithItems) error {
	err := s.validator.ValidateUser(*user)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	err = s.insertUser(ctx, db, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	uw.UserID = user.ID
	uw.OwnerID = user.ID
	err = s.createWorkspace(ctx, db, uw, workspace)
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	for _, ci := range cardsWithItems {
		card := ci.Card
		card.WorkspaceID = workspace.ID
		err = s.createCardWithItems(ctx, db, &card, ci.CardItems)
		if err != nil {
			return fmt.Errorf("failed to create a card: %w", err)
		}
	}

	return nil
}

func (s *Storage) updateCredentials(ctx context.Context, db IConn, ID string, oldIdentity string, newUser domain.User) error {
	oldUser, err := s.findUser(ctx, db, ID)
	if err != nil {
		return fmt.Errorf("faield to fetch a user: %w", err)
	}

	err = s.validator.ValidateUpdateCredentials(validator.ValidateUpdateCredentialsOpts{
		OldUser:     oldUser,
		NewUser:     newUser,
		OldIdentity: oldIdentity,
	})
	if err != nil {
		return fmt.Errorf("failed to validate update password: %w", err)
	}

	query := Builder().Update("users").
		Set("identity", newUser.Identity).
		Set("verifier", Expr("encrypt(?, ?, 'aes-cbc/pad:pkcs')", newUser.Verifier, s.sseKey)).
		Set("srp_salt", newUser.SrpSalt).
		Set("passwd_salt", newUser.PasswdSalt).
		Set("priv_key_enc", newUser.PrivKeyEnc).
		Where("id = ?", ID)

	return s.Exec1(ctx, db, query)
}

func (s *Storage) insertUser(ctx context.Context, db IConn, user *domain.User) error {
	query := Builder().
		Insert("users").
		Columns(
			"ver",
			"readonly",
			"identity",
			"verifier",
			"otp_secret",
			"otp_candidate",
			"srp_salt",
			"passwd_salt",
			"priv_key_enc",
			"pub_key").
		Values(
			user.Ver,
			user.Readonly,
			user.Identity,
			Expr("encrypt(?, ?, 'aes-cbc/pad:pkcs')", user.Verifier, s.sseKey),
			Expr("encrypt(?, ?, 'aes-cbc/pad:pkcs')", user.OtpSecret, s.sseKey),
			Expr("encrypt(?, ?, 'aes-cbc/pad:pkcs')", user.OtpCandidate, s.sseKey),
			user.SrpSalt,
			user.PasswdSalt,
			user.PrivKeyEnc,
			user.PubKey).
		Suffix("RETURNING id, created_at")

	return s.QueryRow(ctx, db, query).Scan(&user.ID, &user.CreatedAt)
}

func (s *Storage) selectUserColumns() SelectBuilder {
	return Builder().
		Select("*").
		Column("decrypt(verifier, ?, 'aes-cbc/pad:pkcs') AS verifier", s.sseKey).
		Column("decrypt(otp_secret, ?, 'aes-cbc/pad:pkcs') AS otp_secret", s.sseKey).
		Column("decrypt(otp_candidate, ?, 'aes-cbc/pad:pkcs') AS otp_candidate", s.sseKey)
}
