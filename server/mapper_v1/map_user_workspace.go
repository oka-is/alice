package mapper_v1

import (
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
)

func MapUserWorkspace(input domain.UserWorkspace) *alice_v1.UserWorkspace {
	return &alice_v1.UserWorkspace{
		Id:          input.ID.String,
		UserId:      input.UserID.String,
		OwnerId:     input.OwnerID.String,
		WorkspaceId: input.WorkspaceID.String,
	}
}

func MapUserWorkspaces(input []domain.UserWorkspace) []*alice_v1.UserWorkspace {
	out := make([]*alice_v1.UserWorkspace, len(input))

	for ix := range input {
		out[ix] = MapUserWorkspace(input[ix])
	}

	return out
}
