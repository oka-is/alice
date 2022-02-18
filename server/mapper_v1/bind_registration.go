package mapper_v1

import (
	"github.com/oka-is/alice/desc/alice_v1"
	"github.com/oka-is/alice/pkg/domain"
)

func BindRegistration(input *alice_v1.RegistrationRequest) (domain.User, domain.UserWorkspace, domain.Workspace, []domain.CardWithItems) {
	u, w := input.GetUser(), input.GetWorkspace()

	user := domain.User{
		Ver:        domain.NewEmptyInt64(1),
		Identity:   domain.NewEmptyString(u.GetIdentity()),
		Verifier:   domain.NewEmptyBytes(u.GetVerifier()),
		SrpSalt:    domain.NewEmptyBytes(u.GetSrpSalt()),
		PasswdSalt: domain.NewEmptyBytes(u.GetPasswdSalt()),
		PrivKeyEnc: domain.NewEmptyBytes(u.GetPrivKeyEnc()),
		PubKey:     domain.NewEmptyBytes(u.GetPubKey()),
	}

	userWorkspace := domain.UserWorkspace{
		AedKeyEnc: domain.NewEmptyBytes(w.GetAedKeyEnc()),
		AedKeyTag: domain.NewEmptyBytes(w.GetAedKeyTag()),
	}

	workspace := domain.Workspace{
		TitleEnc: domain.NewEmptyBytes(w.GetTitleEnc()),
	}

	cards := make([]domain.CardWithItems, len(input.GetCardWithItems()))
	for cx, cw := range input.GetCardWithItems() {
		items := make([]domain.CardItem, len(cw.GetItems()))
		for ix, item := range cw.GetItems() {
			items[ix] = domain.CardItem{
				TitleEnc: domain.NewEmptyBytes(item.GetTitleEnc()),
				BodyEnc:  domain.NewEmptyBytes(item.GetBodyEnc()),
			}
		}

		cards[cx] = domain.CardWithItems{
			CardItems: items,
			Card: domain.Card{
				TitleEnc: domain.NewEmptyBytes(cw.GetTitleEnc()),
				TagsEnc:  domain.NewEmptyByteSlice(cw.GetTagsEnc()),
			},
		}
	}

	return user, userWorkspace, workspace, cards
}
