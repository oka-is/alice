package mapper_v1

import (
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
)

func MapWhoAmI(user domain.User) *alice_v1.WhoAmIResponse {
	return &alice_v1.WhoAmIResponse{
		User: MapPrivUser(user),
	}
}
