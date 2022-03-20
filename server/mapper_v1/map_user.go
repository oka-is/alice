package mapper_v1

import (
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
)

func MapPrivUser(input domain.User) *alice_v1.PrivUser {
	return &alice_v1.PrivUser{
		Id:         input.ID.String,
		Readonly:   input.Readonly.Bool,
		PasswdSalt: input.PasswdSalt.Bytea,
		PrivKeyEnc: input.PrivKeyEnc.Bytea,
		PubKey:     input.PubKey.Bytea,
		OtpEnabled: input.IsOtpEnabled(),
	}
}
