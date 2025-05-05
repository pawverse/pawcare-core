package common

import (
	"context"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/golang-jwt/jwt/v5"
)

func GetStandardClaimsFromContext(ctx context.Context) (*jwt.RegisteredClaims, error) {
	claims, ok := ctx.Value(kitjwt.JWTClaimsContextKey).(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrParsingClaims
	}
	return claims, nil
}

func GetUserIdFromContext(ctx context.Context) (string, error) {
	claims, err := GetStandardClaimsFromContext(ctx)
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}
