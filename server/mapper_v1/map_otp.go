package mapper_v1

import (
	"github.com/wault-pw/alice/desc/alice_v1"
)

func MapOtpIssue(secret, url string) *alice_v1.OtpIssueResponse {
	return &alice_v1.OtpIssueResponse{
		Secret: secret,
		Url:    url,
	}
}
