package storage

import (
	"context"
	"testing"

	"github.com/oka-is/alice/pkg/domain"
)

func mustBuildCardItem(t *testing.T, storage *Storage, input *domain.CardItem) *domain.CardItem {
	out := &domain.CardItem{
		TitleEnc: domain.NewEmptyBytes([]byte("TitleEnc")),
		BodyEnc:  domain.NewEmptyBytes([]byte("BodyEnc")),
		Position: domain.NewNullInt64(0),
	}

	if input.TitleEnc.Valid {
		out.TitleEnc = input.TitleEnc
	}

	if input.BodyEnc.Valid {
		out.BodyEnc = input.BodyEnc
	}

	if input.Position.Valid {
		out.Position = input.Position
	}

	if input.CardID.Valid {
		out.CardID = input.CardID
	} else {
		out.CardID = mustCreateCard(t, storage, &domain.Card{}).ID
	}

	return out
}

func mustCreateCardItem(t *testing.T, storage *Storage, input *domain.CardItem) *domain.CardItem {
	ctx := context.Background()
	output := mustBuildCardItem(t, storage, input)
	if err := storage.upsertCardItem(ctx, storage.db, output); err != nil {
		t.Errorf("cant create factory card: %s", err)
		t.FailNow()
	}
	return output
}
