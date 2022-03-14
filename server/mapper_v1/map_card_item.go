package mapper_v1

import (
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
)

func MapCardItem(input domain.CardItem) *alice_v1.CardItem {
	return &alice_v1.CardItem{
		Id:       input.ID.String,
		CardId:   input.CardID.String,
		Position: input.Position.Int64,
		TitleEnc: input.TitleEnc.Bytea,
		BodyEnc:  input.BodyEnc.Bytea,
		Hidden:   input.Hidden.Bool,
	}
}

func MapCardItems(input []domain.CardItem) []*alice_v1.CardItem {
	out := make([]*alice_v1.CardItem, len(input))

	for ix := range input {
		out[ix] = MapCardItem(input[ix])
	}

	return out
}
