package storage

import (
	"context"
	"fmt"

	"github.com/oka-is/alice/pkg/domain"
	"github.com/oka-is/alice/pkg/validator"
)

func (s *Storage) CreateUser(ctx context.Context, user *domain.User, uw *domain.UserWorkspace,
	workspace *domain.Workspace, cardsWithItems []domain.CardWithItems) error {
	return s.Tx(ctx, nil, func(c context.Context, tx IConn) error {
		return s.createUser(c, tx, user, uw, workspace, cardsWithItems)
	})
}

func (s *Storage) FindUserIdentity(ctx context.Context, identity string) (user domain.User, err error) {
	query := Builder().Select("*").From("users").Where("identity = ?", identity).Limit(1)
	err = s.Get(ctx, s.db, &user, query)
	return
}

func (s *Storage) FindUser(ctx context.Context, ID string) (user domain.User, err error) {
	return s.findUserDB(ctx, s.db, ID)
}

func (s *Storage) TerminateUser(ctx context.Context, identity []byte, userID string) error {
	return s.Tx(ctx, nil, func(c context.Context, tx IConn) error {
		return s.terminateUserDB(c, tx, identity, userID)
	})
}

func (s *Storage) terminateUserDB(ctx context.Context, db IConn, identity []byte, userID string) error {
	user, err := s.findUserDB(ctx, db, userID)
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

func (s *Storage) findUserDB(ctx context.Context, db IConn, ID string) (user domain.User, err error) {
	query := Builder().Select("*").From("users").Where("id = ?", ID).Limit(1)
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

func (s *Storage) insertUser(ctx context.Context, db IConn, user *domain.User) error {
	query := Builder().
		Insert("users").
		Columns(
			"ver",
			"identity",
			"verifier",
			"srp_salt",
			"passwd_salt",
			"priv_key_enc",
			"pub_key").
		Values(
			user.Ver,
			user.Identity,
			user.Verifier,
			user.SrpSalt,
			user.PasswdSalt,
			user.PrivKeyEnc,
			user.PubKey).
		Suffix("RETURNING id, created_at")

	return s.QueryRow(ctx, db, query).Scan(&user.ID, &user.CreatedAt)
}
