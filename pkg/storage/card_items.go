package storage

import (
	"context"

	"github.com/wault-pw/alice/pkg/domain"
)

func (s *Storage) ListCardItems(ctx context.Context, cardID string) (out []domain.CardItem, err error) {
	query := Builder().Select("*").From("card_items").Where("card_id = ?", cardID)
	err = s.Select(ctx, s.db, &out, query)
	return
}

func (s *Storage) upsertCardItem(ctx context.Context, db IConn, item *domain.CardItem) error {
	query := Builder().
		Insert("card_items").
		Columns("card_id", "position", "title_enc", "body_enc", "hidden").
		Values(item.CardID, item.Position, item.TitleEnc, item.BodyEnc, item.Hidden).
		Suffix(`
ON CONFLICT (card_id, position) DO UPDATE SET
title_enc = EXCLUDED.title_enc,
body_enc = EXCLUDED.body_enc,
hidden = EXCLUDED.hidden`).
		Suffix("RETURNING id")

	return s.QueryRow(ctx, db, query).Scan(&item.ID)
}

func (s *Storage) deleteCardItemsPositionedGT(ctx context.Context, db IConn, cardID string, position int) error {
	query := Builder().Delete("card_items").Where("card_id = ?", cardID).Where("position > ?", position)
	return s.Exec1(ctx, db, query)
}

func (s *Storage) cloneCardItems(ctx context.Context, db IConn, newCardID, oldCardID string) error {
	sub := Builder().Select().
		From("card_items").
		Column("? AS card_id", newCardID).
		Column("position").
		Column("title_enc").
		Column("body_enc").
		Column("hidden").
		Where("card_id = ?", oldCardID)

	query := Builder().
		Insert("card_items").
		Columns("card_id", "position", "title_enc", "body_enc", "hidden").Select(sub)

	return s.Exec1(ctx, db, query)
}
