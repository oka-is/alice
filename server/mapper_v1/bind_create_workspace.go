package mapper_v1

import (
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
)

func BindCreateWorkspace(input *alice_v1.CreateWorkspaceRequest) (domain.UserWorkspace, domain.Workspace) {
	userWorkspace := domain.UserWorkspace{
		AedKeyEnc: domain.NewEmptyBytes(input.GetAedKeyEnc()),
	}

	workspace := domain.Workspace{
		TitleEnc: domain.NewEmptyBytes(input.GetTitleEnc()),
	}

	return userWorkspace, workspace
}
