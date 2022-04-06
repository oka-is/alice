package policy

import "github.com/wault-pw/alice/pkg/domain"

type WorkspacePolicy struct {
	uw   domain.UserWorkspace
	user domain.User
}

func (w *WorkspacePolicy) Wrap(user domain.User, uw domain.UserWorkspace) IWorkspacePolicy {
	w.user = user
	w.uw = uw
	return w
}

// CanManageWorkspace only the owner can manage the workspace
func (w *WorkspacePolicy) CanManageWorkspace() error {
	if w.uw.OwnerID.String != w.user.ID.String {
		return ErrDenied
	}

	return nil
}

func (w *WorkspacePolicy) CanSeeWorkspace() error {
	if w.uw.UserID.String != w.user.ID.String {
		return ErrDenied
	}

	return nil
}

func (w *WorkspacePolicy) CanDeleteShare() error {
	// user can not delete self-created user workspace
	if w.user.ID == w.uw.UserID && w.user.ID == w.uw.OwnerID {
		return ErrDenied
	}

	if w.user.ID == w.uw.UserID || w.user.ID == w.uw.OwnerID {
		return nil
	}

	return ErrDenied
}

func (w *WorkspacePolicy) CanSeeCard(card domain.Card) error {
	if err := w.CanSeeWorkspace(); err != nil {
		return err
	}

	if card.WorkspaceID.String != w.uw.WorkspaceID.String {
		return ErrDenied
	}

	return nil
}

func (w *WorkspacePolicy) CanManageCard(card domain.Card) error {
	if err := w.CanSeeWorkspace(); err != nil {
		return err
	}

	if card.WorkspaceID.String != w.uw.WorkspaceID.String {
		return ErrDenied
	}

	return nil
}
