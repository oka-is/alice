package jwt

import (
	j "github.com/golang-jwt/jwt"
)

type Algo j.SigningMethod

var (
	HS256 = j.SigningMethodHS256
	HS384 = j.SigningMethodHS384
	HS512 = j.SigningMethodHS512
)

const (
	jti = "jti"
	aud = "aud"
	sub = "sub"
	iss = "iss"
	exp = "exp"
)
