package backup

import (
	"fmt"

	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/server/mapper_v1"
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

		cards, err := b.store.ListCardsByWorkspace(b.ctx, workspaces[ix].WorkspaceID.String)
		if err != nil {
			return fmt.Errorf("filed to list cards: %w", err)
		}

		for cid := range cards {
			if err = ListCard(b, cards[cid]); err != nil {
				return fmt.Errorf("failed to dump a card: %w", err)
			}

			if err = ListCardItems(b, cards[cid]); err != nil {
				return fmt.Errorf("failed to dump card items: %w", err)
			}
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

func ListCard(b *Backup, card domain.Card) error {
	_, err := b.writer.write(MarkerCard, mapper_v1.MapCard(card))
	if err != nil {
		return err
	}

	b.f.Flush()
	return nil
}

func ListCardItems(b *Backup, card domain.Card) error {
	items, err := b.store.ListCardItems(b.ctx, card.ID.String)
	if err != nil {
		return fmt.Errorf("filed to list card items: %w", err)
	}

	for ix := range items {
		_, err = b.writer.write(MarkerCardItem, mapper_v1.MapCardItem(items[ix]))
		if err != nil {
			return err
		}
	}

	b.f.Flush()
	return nil
}
