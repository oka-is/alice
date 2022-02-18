package jwt

import "time"

type IOts interface {
	Marshall() (string, error)
	Unmarshall(input string) (string, error)
	SetJti(jti string) IOts
	SetExp(exp time.Time) IOts
}
