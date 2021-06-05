package service

import (
	"context"
	"errors"

	"github.com/StevenRojas/goaccess-api/pkg/utils"
	pkgServ "github.com/StevenRojas/goaccess/pkg/service"
)

// AppService app service
type AppService interface {
	// Login validate email account and generate a JTW token
	Login(ctx context.Context, email string) (string, error)
}

type appService struct {
	authService pkgServ.AuthenticationService
	jh          utils.JwtHandler
}

// NewAppService return a new app service instance
func NewAppService(authService pkgServ.AuthenticationService, jsonHandler utils.JwtHandler) *appService {
	return &appService{
		authService: authService,
		jh:          jsonHandler,
	}
}

// Login validate email account and generate a JTW token
func (a *appService) Login(ctx context.Context, email string) (string, error) {
	user, err := a.authService.Login(ctx, email)
	if err != nil {
		return "", err
	}
	if !user.User.IsAdmin {
		return "", errors.New("you don't have permissions to access this module")
	}
	token, err := a.jh.CreateToken(user.User.ID)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}
