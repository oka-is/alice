package mapper_v1

import (
	"github.com/oka-is/alice/desc/alice_v1"
	"github.com/oka-is/alice/pkg/domain"
)

func BindUpsertCard(req *alice_v1.UpsertCardRequest) (domain.Card, []domain.CardItem) {
	card := BindCard(req.GetCard())
	items := BindCardItems(req.GetCardItems())
	return card, items
}

func BindCard(req *alice_v1.Card) domain.Card {
	return domain.Card{
		WorkspaceID: domain.NewNullString(),
		Archived:    domain.NewEmptyBool(req.GetArchived()),
		TitleEnc:    domain.NewEmptyBytes(req.GetTitleEnc()),
		TagsEnc:     domain.NewEmptyByteSlice(req.GetTagsEnc()),
	}
}

func BindCardItems(req []*alice_v1.CardItem) []domain.CardItem {
	out := make([]domain.CardItem, len(req))
	for ix, item := range req {
		out[ix] = BindCardItem(item)
	}
	return out
}

func BindCardItem(req *alice_v1.CardItem) domain.CardItem {
	return domain.CardItem{
		CardID:   domain.NewNullString(),
		TitleEnc: domain.NewEmptyBytes(req.GetTitleEnc()),
		BodyEnc:  domain.NewEmptyBytes(req.GetBodyEnc()),
		Hidden:   domain.NewEmptyBool(req.GetHidden()),
	}
}
