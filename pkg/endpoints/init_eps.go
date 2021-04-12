package endpoints

import (
	"context"
	"errors"
	"os"

	"github.com/StevenRojas/goaccess-api/pkg/codec"
	e "github.com/StevenRojas/goaccess-api/pkg/errors"
	"github.com/StevenRojas/goaccess/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

type InitEndpoints struct {
	ForceReset endpoint.Endpoint
}

func MakeInitEndpoints(
	s service.InitializationService,
	middlewares []endpoint.Middleware) InitEndpoints {
	return InitEndpoints{
		ForceReset: wrapMiddlewares(makeForceReset(s), middlewares),
	}
}

func makeForceReset(s service.InitializationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reset, ok := os.LookupEnv("ALLOW_RESET")
		if !ok || reset != "true" {
			return nil, e.HTTPForbidden(errors.New("Reset the DB is not allowed"))
		}
		err := s.Init(ctx, false)
		if err != nil {
			return nil, e.HTTPConflict("Unable to reset the database", err)
		}
		return &codec.EmptyResponse{}, nil
	}
}
