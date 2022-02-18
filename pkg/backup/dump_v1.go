package backup

import (
	"fmt"

	"github.com/oka-is/alice/pkg/domain"
	"github.com/oka-is/alice/server/mapper_v1"
)

func (b *Backup) DumpV1(userID string) error {
	workspaces, err := b.store.ListUserWithWorkspaces(b.ctx, userID)
	if err != nil {
		return fmt.Errorf("filed to list user workspaces: %w", err)
	}

	if err = Whoami(b, userID); err != nil {
		return fmt.Errorf("failed to dump a whoami: %w", err)
	}

	if err = ListWorkspacesResponse(b, workspaces); err != nil {
		return fmt.Errorf("failed to dump a workspaces: %w", err)
	}

	if err = ListCardsResponses(b, workspaces); err != nil {
		return fmt.Errorf("failed to dump a workspaces: %w", err)
	}

	return nil
}

func Whoami(b *Backup, userID string) error {
	user, err := b.store.FindUser(b.ctx, userID)
	if err != nil {
		return fmt.Errorf("filed to find a user: %w", err)
	}

	if _, err = b.writer.write(
		MarkerWhoAmI,
		mapper_v1.MapWhoAmI(user)); err != nil {
		return err
	}

	b.f.Flush()
	return nil
}

func ListWorkspacesResponse(b *Backup, workspaces []domain.UserWithWorkspace) error {
	if _, err := b.writer.write(
		MarkerListWorkspacesResponse,
		mapper_v1.MapListUserWorkspaceResponse(workspaces)); err != nil {
		return err
	}

	b.f.Flush()
	return nil
}

func ListCardsResponses(b *Backup, workspaces []domain.UserWithWorkspace) error {
	for ix := range workspaces {
		cards, err := b.store.ListCardsByWorkspace(b.ctx, workspaces[ix].WorkspaceID.String)
		if err != nil {
			return fmt.Errorf("filed to list cards: %w", err)
		}

		if _, err = b.writer.write(
			MarkerListCardsResponse,
			mapper_v1.MapListCardsResponse(cards)); err != nil {
			return err
		}

		b.f.Flush()
	}

	return nil
}
