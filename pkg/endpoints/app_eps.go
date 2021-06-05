package endpoints

import (
	"context"
	"errors"

	"github.com/StevenRojas/goaccess-api/pkg/codec"
	e "github.com/StevenRojas/goaccess-api/pkg/errors"
	appServ "github.com/StevenRojas/goaccess-api/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

type AppEndpoints struct {
	Login endpoint.Endpoint
}

func MakeAppEndpoints(
	s appServ.AppService,
	middlewares []endpoint.Middleware) AppEndpoints {
	return AppEndpoints{
		Login: wrapMiddlewares(makeLogin(s), middlewares),
	}
}

func makeLogin(s appServ.AppService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		email, ok := request.(string)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to string"))
		}
		token, err := s.Login(ctx, email)
		if err != nil {
			return nil, e.HTTPForbidden(err)
		}
		return &codec.TokenResponse{
			Token: token,
		}, nil
	}
}
