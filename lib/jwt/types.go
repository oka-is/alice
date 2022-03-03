package jwt

import "time"

//go:generate mockgen -destination ../jwt_mock/ots_mock.go -source types.go -package jwt_mock -mock_names IOts=MockOts
type IOts interface {
	Marshall() (string, error)
	Unmarshall(input string) (string, error)
	SetJti(jti string) IOts
	SetExp(exp time.Time) IOts
}
