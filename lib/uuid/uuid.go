package uuid

import (
	"strings"

	u "github.com/google/uuid"
)

const Empty = "00000000-0000-0000-0000-000000000000"

func NewV4() string {
	return u.New().String()
}

func Safe(input string) string {
	uuid, err := u.Parse(strings.ReplaceAll(input, " ", ""))
	if err != nil {
		return Empty
	}

	return uuid.String()
}
