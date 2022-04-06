package domain

import "github.com/wault-pw/alice/lib/uuid"

func NewUUID() string {
	return uuid.NewV4()
}
