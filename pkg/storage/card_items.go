package storage

import (
	"context"

	"github.com/oka-is/alice/pkg/domain"
)

func (s *Storage) ListCardItems(ctx context.Context, cardID string) (out []domain.CardItem, err error) {
	query := Builder().Select("*").From("card_items").Where("card_id = ?", cardID)
	err = s.Select(ctx, s.db, &out, query)
	return
}

func (s *Storage) insertCardItem(ctx context.Context, db IConn, item domain.CardItem) error {
	query := Builder().
		Insert("card_items").
		Columns("card_id", "title_enc", "body_enc", "hidden").
		Values(item.CardID, item.TitleEnc, item.BodyEnc, item.Hidden).
		Suffix("RETURNING id")

	return s.QueryRow(ctx, db, query).Scan(&item.ID)
}

func (s *Storage) cloneCardItems(ctx context.Context, db IConn, newCardID, oldCardID string) error {
	sub := Builder().Select().
		From("card_items").
		Column("? AS card_id", newCardID).
		Column("title_enc").
		Column("body_enc").
		Column("hidden").
		Where("card_id = ?", oldCardID)

	query := Builder().
		Insert("card_items").
		Columns("card_id", "title_enc", "body_enc", "hidden").Select(sub)
	_, err := s.Exec(ctx, db, query)
	return err
}
