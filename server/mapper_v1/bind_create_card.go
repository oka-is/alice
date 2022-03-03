package mapper_v1

import (
	"github.com/oka-is/alice/desc/alice_v1"
	"github.com/oka-is/alice/pkg/domain"
)

func BindCreateCard(input *alice_v1.CreateCardRequest) (domain.Card, []domain.CardItem) {
	card := BindCard(input.GetCard())
	items := BindCardItems(input.GetCardItems())
	return card, items
}

func BindCard(input *alice_v1.Card) domain.Card {
	return domain.Card{
		WorkspaceID: domain.NewNullString(),
		TitleEnc:    domain.NewEmptyBytes(input.GetTitleEnc()),
	}
}

func BindCardItems(input []*alice_v1.CardItem) []domain.CardItem {
	out := make([]domain.CardItem, len(input))
	for ix, item := range input {
		out[ix] = BindCardItem(item)
	}
	return out
}

func BindCardItem(input *alice_v1.CardItem) domain.CardItem {
	return domain.CardItem{
		CardID:   domain.NewNullString(),
		TitleEnc: domain.NewEmptyBytes(input.GetTitleEnc()),
		BodyEnc:  domain.NewEmptyBytes(input.GetBodyEnc()),
		Hidden:   domain.NewEmptyBool(input.GetHidden()),
	}
}
