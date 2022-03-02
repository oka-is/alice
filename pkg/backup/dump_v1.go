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

	for ix := range workspaces {
		if err = ListWorkspace(b, workspaces[ix]); err != nil {
			return fmt.Errorf("failed to dump a workspaces: %w", err)
		}

		if err = ListCards(b, workspaces[ix]); err != nil {
			return fmt.Errorf("failed to dump a cards: %w", err)
		}
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

func ListWorkspace(b *Backup, workspace domain.UserWithWorkspace) error {
	if _, err := b.writer.write(MarkerWorkspace,
		mapper_v1.MapUserWithWorkspace(workspace)); err != nil {
		return err
	}

	b.f.Flush()
	return nil
}

func ListCards(b *Backup, workspace domain.UserWithWorkspace) error {
	cards, err := b.store.ListCardsByWorkspace(b.ctx, workspace.WorkspaceID.String)
	if err != nil {
		return fmt.Errorf("filed to list cards: %w", err)
	}

	for ix := range cards {
		_, err = b.writer.write(MarkerCard, mapper_v1.MapCard(cards[ix]))
		if err != nil {
			return err
		}
	}

	b.f.Flush()
	return nil
}
