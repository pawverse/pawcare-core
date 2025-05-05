package common

import (
	"context"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func JWTKeyFactory(viper viper.Viper) func(token *jwt.Token) (any, error) {
	return func(token *jwt.Token) (any, error) {
		return []byte(viper.GetString(JWTSecretKey)), nil
	}
}

func RegisteredClaimsFactory() jwt.Claims {
	return &jwt.RegisteredClaims{}
}

type ClaimsFactory func() jwt.Claims

// NewParser creates a new JWT parsing middleware, specifying a
// jwt.Keyfunc interface, the signing method and the claims type to be used. NewParser
// adds the resulting claims to endpoint context or returns error on invalid token.
// Particularly useful for servers.
func NewParser(keyFunc jwt.Keyfunc, method jwt.SigningMethod, newClaims ClaimsFactory) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request any) (response any, err error) {
			// tokenString is stored in the context from the transport handlers.
			tokenString, ok := ctx.Value(kitjwt.JWTContextKey).(string)
			if !ok {
				return nil, kitjwt.ErrTokenContextMissing
			}

			// Parse takes the token string and a function for looking up the
			// key. The latter is especially useful if you use multiple keys
			// for your application.  The standard is to use 'kid' in the head
			// of the token to identify which key to use, but the parsed token
			// (head and claims) is provided to the callback, providing
			// flexibility.
			token, err := jwt.ParseWithClaims(tokenString, newClaims(), func(token *jwt.Token) (any, error) {
				// Don't forget to validate the alg is what you expect:
				if token.Method != method {
					return nil, kitjwt.ErrUnexpectedSigningMethod
				}

				return keyFunc(token)
			})
			if err != nil {
				switch err {
				case jwt.ErrTokenExpired:
					return nil, kitjwt.ErrTokenExpired
				case jwt.ErrTokenMalformed:
					return nil, kitjwt.ErrTokenMalformed
				case jwt.ErrTokenNotValidYet:
					return nil, kitjwt.ErrTokenNotActive
				default:
					return nil, err
				}
			}

			if !token.Valid {
				return nil, kitjwt.ErrTokenInvalid
			}

			ctx = context.WithValue(ctx, kitjwt.JWTClaimsContextKey, token.Claims)

			return next(ctx, request)
		}
	}
}
