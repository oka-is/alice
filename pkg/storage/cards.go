package storage

import (
	"context"
	"fmt"

	"github.com/oka-is/alice/pkg/domain"
)

func (s *Storage) CreateCardWithItems(ctx context.Context, card *domain.Card, items []domain.CardItem) error {
	return s.Tx(ctx, nil, func(c context.Context, tx IConn) error {
		return s.createCardWithItems(c, tx, card, items)
	})
}

func (s *Storage) UpdateCardWithItems(ctx context.Context, card *domain.Card, items []domain.CardItem) error {
	return s.Tx(ctx, nil, func(c context.Context, tx IConn) error {
		return s.updateCardWithItems(c, tx, card, items)
	})
}

func (s *Storage) CloneCard(ctx context.Context, oldCardID string, titleEnc []byte) (out domain.Card, err error) {
	err = s.Tx(ctx, nil, func(c context.Context, tx IConn) error {
		out, err = s.cloneCard(c, tx, oldCardID, titleEnc)
		return err
	})

	return
}

func (s *Storage) ArchiveCard(ctx context.Context, ID string) (archived bool, err error) {
	query := Builder().
		Update("cards").
		Set("archived", Expr("NOT COALESCE(archived, FALSE)")).
		Where("id = ?", ID).
		Suffix("RETURNING archived")

	err = s.Get(ctx, s.db, &archived, query)
	return
}

func (s *Storage) ListCardsByWorkspace(ctx context.Context, workspaceID string) (out []domain.Card, err error) {
	query := Builder().Select("*").From("cards").Where("workspace_id = ?", workspaceID)
	err = s.Select(ctx, s.db, &out, query)
	return
}

func (s *Storage) DeleteCard(ctx context.Context, cardID string) error {
	query := Builder().Delete("cards").Where("id = ?", cardID)
	_, err := s.Exec(ctx, s.db, query)
	return err
}

func (s *Storage) updateCardWithItems(ctx context.Context, db IConn, card *domain.Card, items []domain.CardItem) error {
	if err := s.updateCard(ctx, db, card); err != nil {
		return fmt.Errorf("failed to update card: %w", err)
	}

	for ix := range items {
		item := &items[ix]
		item.CardID = card.ID
		item.Position = domain.NewNullInt64(int64(ix))
		if err := s.upsertCardItem(ctx, db, item); err != nil {
			return fmt.Errorf("failed to update card item: %w", err)
		}
	}

	if err := s.deleteCardItemsPositionedGT(ctx, db, card.ID.String, len(items)-1); err != nil {
		return fmt.Errorf("failed to delete positioned items: %w", err)
	}

	return nil
}

func (s *Storage) createCardWithItems(ctx context.Context, db IConn, card *domain.Card, items []domain.CardItem) error {
	if err := s.insertCard(ctx, db, card); err != nil {
		return fmt.Errorf("failed to create card: %w", err)
	}

	for ix := range items {
		item := &items[ix]
		item.CardID = card.ID
		item.Position = domain.NewNullInt64(int64(ix))
		if err := s.upsertCardItem(ctx, db, item); err != nil {
			return fmt.Errorf("failed to update card item: %w", err)
		}
	}

	return nil
}

func (s *Storage) cloneCard(ctx context.Context, db IConn, oldCardID string, titleEnc []byte) (out domain.Card, err error) {
	sub := Builder().Select().
		From("cards").
		Column("? AS title_enc", titleEnc).
		Column("workspace_id").
		Column("archived").
		Column("tags_enc").
		Where("id = ?", oldCardID).
		Suffix("RETURNING id")

	query := Builder().
		Insert("cards").
		Columns("title_enc", "workspace_id", "archived", "tags_enc").
		Select(sub)

	var newID string
	err = s.QueryRow(ctx, db, query).Scan(&newID)
	if err != nil {
		return out, fmt.Errorf("failed to clone card items: %w", err)
	}

	err = s.cloneCardItems(ctx, db, newID, oldCardID)
	if err != nil {
		return out, fmt.Errorf("failed to clone card items: %w", err)
	}

	return s.findCard(ctx, db, newID)
}

func (s *Storage) findCard(ctx context.Context, db IConn, ID string) (out domain.Card, err error) {
	query := Builder().Select("*").From("cards").Where("id = ?", ID).Limit(1)
	err = s.Get(ctx, db, &out, query)
	return
}

func (s *Storage) insertCard(ctx context.Context, db IConn, card *domain.Card) error {
	query := Builder().
		Insert("cards").
		Columns("workspace_id", "archived", "title_enc", "tags_enc").
		Values(card.WorkspaceID, card.Archived, card.TitleEnc, card.TagsEnc).
		Suffix("RETURNING id, created_at")

	return s.QueryRow(ctx, db, query).Scan(&card.ID, &card.CreatedAt)
}

func (s *Storage) updateCard(ctx context.Context, db IConn, card *domain.Card) error {
	query := Builder().
		Update("cards").
		Set("archived", card.Archived).
		Set("title_enc", card.TitleEnc).
		Set("tags_enc", card.TagsEnc).
		Where("id = ?", card.ID.String)

	return s.Exec1(ctx, db, query)
}
