package storage

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/oka-is/alice/pkg/validator"
)

var (
	pgDSN = os.Getenv("PG_DSN")
	_db   *sqlx.DB
)

func TestMain(m *testing.M) {
	db, err := Connect(pgDSN)
	if err != nil {
		panic(err)
	}

	_db = db
	defer _db.Close()
	os.Exit(m.Run())
}

func MustNewStore(t *testing.T) (*Storage, *SavepointDB) {
	return NewSavepointStorage(_db, validator.NewNoOpt(), []byte{1, 2, 3, 4, 5, 6, 7, 8, 9})
}
