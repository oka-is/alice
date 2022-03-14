package mapper_v1

import (
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
)

func MapListUserWorkspaceResponse(items []domain.UserWithWorkspace) *alice_v1.ListWorkspacesResponse {
	return &alice_v1.ListWorkspacesResponse{
		Items: MapUserWithWorkspaces(items),
	}
}

func MapUserWithWorkspaces(input []domain.UserWithWorkspace) []*alice_v1.UserWithWorkspace {
	out := make([]*alice_v1.UserWithWorkspace, len(input))
	for ix, workspace := range input {
		out[ix] = MapUserWithWorkspace(workspace)
	}
	return out
}

func MapUserWithWorkspace(input domain.UserWithWorkspace) *alice_v1.UserWithWorkspace {
	return &alice_v1.UserWithWorkspace{
		Id:          input.WorkspaceID.String,
		UserId:      input.UserID.String,
		OwnerId:     input.OwnerID.String,
		OwnerPubKey: input.OwnerPubKey.Bytea,
		WorkspaceId: input.WorkspaceID.String,
		AedKeyEnc:   input.AedKeyEnc.Bytea,
		TitleEnc:    input.TitleEnc.Bytea,
		CreatedAt:   input.WorkspaceCreatedAt.Time.String(),
	}
}
