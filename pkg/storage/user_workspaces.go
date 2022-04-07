package storage

import (
	"context"
	"fmt"

	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/pkg/validator"
)

func (s *Storage) FindUserWorkspace(ctx context.Context, ID string) (out domain.UserWorkspace, err error) {
	return s.findUserWorkspace(ctx, s.db, ID)
}

func (s *Storage) FindUserWorkspaceLink(ctx context.Context, userID, workspaceID string) (out domain.UserWorkspace, err error) {
	query := Builder().Select("*").From("user_workspaces").Where("user_id =? AND workspace_id = ?", userID, workspaceID).Limit(1)
	err = s.Get(ctx, s.db, &out, query)
	return
}

func (s *Storage) ListSharedUserWorkspaces(ctx context.Context, workspaceID, ownerID string) (out []domain.UserWorkspace, err error) {
	query := Builder().Select("*").From("user_workspaces").Where("workspace_id = ? AND user_id <> ?", workspaceID, ownerID)
	err = s.Select(ctx, s.db, &out, query)
	return
}

func (s *Storage) DeleteUserWorkspace(ctx context.Context, ID string) error {
	query := Builder().Delete("user_workspaces").Where("id = ?", ID)
	return s.Exec1(ctx, s.db, query)
}

func (s *Storage) ShareUserWorkspace(ctx context.Context, uw *domain.UserWorkspace) error {
	return s.Tx(ctx, nil, func(c context.Context, tx IConn) error {
		return s.shareUserWorkspace(c, tx, uw)
	})
}

func (s *Storage) shareUserWorkspace(ctx context.Context, db IConn, uw *domain.UserWorkspace) error {
	shareExists, err := s.isUserWorkspaceShared(ctx, db, uw.UserID.String, uw.WorkspaceID.String)
	if err != nil {
		return fmt.Errorf("failed to detect user workspace existance: %w", err)
	}

	_, userExists, err := s.findOrBuildUser(ctx, db, uw.UserID.String)
	if err != nil {
		return fmt.Errorf("failed to find a user: %w", err)
	}

	err = s.validator.ValidateShareUserWorkspace(validator.ValidateShareUserWorkspaceOpts{
		UserWorkspace: *uw,
		NotShared:     !shareExists,
		UserExists:    userExists,
	})
	if err != nil {
		return fmt.Errorf("failed to validate share: %w", err)
	}

	err = s.insertUserWorkspace(ctx, db, uw)
	if err != nil {
		return fmt.Errorf("failed to insert user workspace: %w", err)
	}

	return nil
}

func (s *Storage) findUserWorkspace(ctx context.Context, db IConn, ID string) (out domain.UserWorkspace, err error) {
	query := Builder().Select("*").From("user_workspaces").Where("id = ?", ID).Limit(1)
	err = s.Get(ctx, db, &out, query)
	return
}

func (s *Storage) isUserWorkspaceShared(ctx context.Context, db IConn, userID, workspaceID string) (exists bool, err error) {
	query := Builder().
		Select("1").
		From("user_workspaces").
		Where("user_id =? AND workspace_id = ?", userID, workspaceID).
		Limit(1).
		Prefix("SELECT EXISTS(").
		Suffix(")")

	err = s.Get(ctx, db, &exists, query)
	return
}

func (s *Storage) insertUserWorkspace(ctx context.Context, conn IConn, userWorkspace *domain.UserWorkspace) error {
	query := Builder().
		Insert("user_workspaces").
		Columns(
			"user_id",
			"owner_id",
			"workspace_id",
			"aed_key_enc").
		Values(
			userWorkspace.UserID,
			userWorkspace.OwnerID,
			userWorkspace.WorkspaceID,
			userWorkspace.AedKeyEnc).
		Suffix("RETURNING id, created_at")

	return s.QueryRow(ctx, conn, query).Scan(&userWorkspace.ID, &userWorkspace.CreatedAt)
}
