package null

import "database/sql/driver"

type Bytea struct {
	Bytea []byte
	Valid bool
}

// Scan implements the Scanner interface.
func (b *Bytea) Scan(value interface{}) error {
	if value == nil {
		b.Bytea, b.Valid = []byte{}, false
		return nil
	}
	b.Valid = true
	b.Bytea = value.([]byte)
	return nil
}

// Value implements the driver Valuer interface.
func (b Bytea) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bytea, nil
}
