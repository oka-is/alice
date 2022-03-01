package storage

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/oka-is/alice/lib/jwt"
	"github.com/oka-is/alice/pkg/domain"
)

type (
	Result = sql.Result
	Tx     = sqlx.Tx
	Row    = sql.Row
	TxOpts = sql.TxOptions
	TxFunc func(ctx context.Context, tx *Tx) error
)

//go:generate mockgen -destination ../storage_mock/store_mock.go -source types.go -package storage_mock -mock_names IStore=MockStore
type IStore interface {
	// Common

	Ping(ctx context.Context) error

	// Operations about sessions

	IssueSession(ctx context.Context, opts jwt.IOts) (*domain.Session, string, error)
	RetrieveSession(ctx context.Context, opts jwt.IOts, token string) (session domain.Session, err error)
	NominateSession(ctx context.Context, jti string) error
	CandidateSession(ctx context.Context, jti, candidateID string, srp []byte) error
	DeleteSession(ctx context.Context, jti string) error

	// Operations about users

	CreateUser(ctx context.Context, user *domain.User, uw *domain.UserWorkspace, workspace *domain.Workspace, cards []domain.CardWithItems) error
	FindUserIdentity(ctx context.Context, identity string) (user domain.User, err error)
	FindUser(ctx context.Context, ID string) (user domain.User, err error)

	// Operations about cards & items

	CreateCardWithItems(ctx context.Context, card *domain.Card, items []domain.CardItem) error
	ListCardsByWorkspace(ctx context.Context, workspaceID string) (out []domain.Card, err error)
	ListCardItems(ctx context.Context, cardID string) (out []domain.CardItem, err error)
	DeleteCard(ctx context.Context, cardID string) error
	CloneCard(ctx context.Context, card *domain.Card, oldCardID string) error
	ArchiveCard(ctx context.Context, ID string) (archived bool, err error)

	// Operations about workspaces

	ListUserWithWorkspaces(ctx context.Context, userID string) (out []domain.UserWithWorkspace, err error)
	CreateWorkspace(ctx context.Context, uw *domain.UserWorkspace, workspace *domain.Workspace) error
	FindUserWithWorkspace(ctx context.Context, ID string) (out domain.UserWithWorkspace, err error)
	DeleteWorkspace(ctx context.Context, ID string) error
}

type IBuilder interface {
	ToSql() (string, []interface{}, error)
}

type IConn interface {
	SelectContext(ctx context.Context, des interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row
}
