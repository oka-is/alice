package jwt

import (
	"fmt"
	"time"

	j "github.com/golang-jwt/jwt"
)

type Opts struct {
	// Aud is audience of the token, example is "WEB" or "MOBILE"
	Aud string
	// Sub is subject of the token, example is "LOGIN"
	Sub string
	// Iss is issuer, example is "Apple ink."
	Iss string
	// Jti is a token ID
	Jti  string
	Algo Algo
	Key  []byte
	// Exp is an expiration time
	Exp time.Time
}

func (o Opts) SetJti(jti string) IOts {
	o.Jti = jti
	return o
}

func (o Opts) SetExp(exp time.Time) IOts {
	o.Exp = exp
	return o
}

// Marshall constructs a token from passed options
func (o Opts) Marshall() (string, error) {
	if err := o.validate(); err != nil {
		return "", fmt.Errorf("validation failed: %w", err)
	}

	return j.NewWithClaims(o.Algo, o.claims()).SignedString(o.Key)
}

// Unmarshall parses JWT token string and returns jti
func (o Opts) Unmarshall(input string) (string, error) {
	token, err := j.Parse(input, func(token *j.Token) (interface{}, error) {
		if token.Method.Alg() != o.Algo.Alg() {
			return nil, fmt.Errorf("%w: %s", ErrAlgoMismatch, token.Header["alg"])
		}

		return o.Key, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(j.MapClaims); ok && token.Valid {
		return claims["jti"].(string), nil
	}

	return "", ErrInvalid
}

func (o Opts) validate() error {
	if o.Sub == "" {
		return fmt.Errorf("%w: sub", ErrRequired)
	}

	if o.Aud == "" {
		return fmt.Errorf("%w: aud", ErrRequired)
	}

	if o.Iss == "" {
		return fmt.Errorf("%w: iss", ErrRequired)
	}

	if o.Algo == nil {
		return fmt.Errorf("%w: signature", ErrRequired)
	}

	if len(o.Key) == 0 {
		return fmt.Errorf("%w: key", ErrRequired)
	}

	return nil
}

func (o *Opts) claims() j.MapClaims {
	return j.MapClaims{
		jti: o.Jti,
		aud: o.Aud,
		sub: o.Sub,
		iss: o.Iss,
		exp: o.Exp.Unix(),
	}
}
