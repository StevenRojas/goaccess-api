package middlewares

import (
	"context"
	"errors"
	"net/http"

	e "github.com/StevenRojas/goaccess-api/pkg/errors"
	"github.com/dgrijalva/jwt-go"
	kitJWT "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	gokitHTTP "github.com/go-kit/kit/transport/http"

	log "github.com/StevenRojas/goaccess/pkg/configuration"
)

// JWTCheck check if JWT is valid
func JWTCheck(logger log.LoggerWrapper) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			claimsValue := ctx.Value(kitJWT.JWTClaimsContextKey)
			if claimsValue == nil {
				err := errors.New("invalid JWT claims")
				logger.Error(err.Error())
				return nil, e.HTTPBadRequest(err)
			}
			isErr, ok := claimsValue.(error)
			if ok {
				logger.Error(isErr.Error())
				return nil, e.HTTPUnauthorized(isErr)
			}

			return next(ctx, request)
		}
	}
}

// DecodeJWT decode JWT into claims
func DecodeJWT(signingMethod jwt.SigningMethod, secret string, logger log.LoggerWrapper) (gokitHTTP.RequestFunc, error) {
	return func(ctx context.Context, r *http.Request) context.Context {
		tokenString, ok := ctx.Value(kitJWT.JWTTokenContextKey).(string)
		if !ok {
			return context.WithValue(ctx, kitJWT.JWTClaimsContextKey, errors.New("unable to parse JWT token"))
		}

		token, err := GetTokenClaims(tokenString, secret)
		if err != nil {
			return context.WithValue(ctx, kitJWT.JWTClaimsContextKey, err)
		}

		return context.WithValue(ctx, kitJWT.JWTClaimsContextKey, token)
	}, nil
}

func GetTokenClaims(token string, secret string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("wrong JWT signed method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok && !parsed.Valid {
		return nil, err
	}
	return claims, nil
}
