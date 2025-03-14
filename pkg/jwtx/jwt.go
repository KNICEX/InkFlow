package jwtx

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims[T any] struct {
	jwt.RegisteredClaims
	Data T
}

func Generate[T any](data T, duration time.Duration, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims[T]{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		Data: data,
	})
	return token.SignedString([]byte(secret))
}

func Parse[T any](tokenString string, secret string) (T, error) {
	var claims Claims[T]
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return claims.Data, err
	}
	if !token.Valid {
		return claims.Data, jwt.ErrInvalidKey
	}
	return claims.Data, nil
}
