package mapper_v1

import (
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
)

func BindUpdateCredentials(req *alice_v1.UpdateCredentialsRequest) (string, domain.User) {
	user := domain.User{
		Identity:   domain.NewEmptyString(req.GetNewIdentity()),
		Verifier:   domain.NewEmptyBytes(req.GetVerifier()),
		SrpSalt:    domain.NewEmptyBytes(req.GetSrpSalt()),
		PasswdSalt: domain.NewEmptyBytes(req.GetPasswdSalt()),
		PrivKeyEnc: domain.NewEmptyBytes(req.GetPrivKeyEnc()),
	}

	return req.GetOldIdentity(), user
}
