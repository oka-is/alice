package storage

import (
	"context"
	"fmt"

	"github.com/oka-is/alice/pkg/domain"
)

func (s *Storage) CreateCardWithItems(ctx context.Context, card *domain.Card, items []domain.CardItem) error {
	return s.Tx(ctx, nil, func(c context.Context, tx *Tx) error {
		return s.createCardWithItems(c, tx, card, items)
	})
}

func (s *Storage) CloneCard(ctx context.Context, card *domain.Card, oldCardID string) error {
	return s.Tx(ctx, nil, func(c context.Context, tx *Tx) error {
		return s.cloneCard(c, tx, card, oldCardID)
	})
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

func (s *Storage) createCardWithItems(ctx context.Context, db IConn, card *domain.Card, items []domain.CardItem) error {
	if err := s.insertCard(ctx, db, card); err != nil {
		return fmt.Errorf("failed to create card: %w", err)
	}

	for _, item := range items {
		item.CardID = card.ID
		if err := s.insertCardItem(ctx, db, item); err != nil {
			return fmt.Errorf("failed to create card item: %w", err)
		}
	}

	return nil
}

func (s *Storage) cloneCard(ctx context.Context, db IConn, card *domain.Card, oldCardID string) error {
	err := s.insertCard(ctx, db, card)
	if err != nil {
		return fmt.Errorf("failed to insert card: %w", err)
	}

	err = s.cloneCardItems(ctx, db, card.ID.String, oldCardID)
	if err != nil {
		return fmt.Errorf("failed to clone card items: %w", err)
	}

	return nil
}

func (s *Storage) insertCard(ctx context.Context, db IConn, card *domain.Card) error {
	query := Builder().
		Insert("cards").
		Columns("workspace_id", "title_enc", "tags_enc").
		Values(card.WorkspaceID, card.TitleEnc, card.TagsEnc).
		Suffix("RETURNING id, created_at")

	return s.QueryRow(ctx, db, query).Scan(&card.ID, &card.CreatedAt)
}

func (s *Storage) insertCardItem(ctx context.Context, db IConn, item domain.CardItem) error {
	query := Builder().
		Insert("card_items").
		Columns("card_id", "title_enc", "body_enc").
		Values(item.CardID, item.TitleEnc, item.BodyEnc).
		Suffix("RETURNING id")

	return s.QueryRow(ctx, db, query).Scan(&item.ID)
}
