package server

import (
	"context"
	"github.com/StevenRojas/goaccess-api/pkg/pb"
	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServer struct {
	authenticationService service.AuthenticationService
}

// NewUserServer return a User Server instance
func NewUserServer(authenticationService service.AuthenticationService) pb.UsersServer {
	return &userServer{
		authenticationService,
	}
}

// RegisterUser register a new user
func (u userServer) RegisterUser(ctx context.Context, request *pb.UserRequest) (*pb.EmptyResponse, error) {
	user := &entities.User{
		ID:      request.Id,
		Email:   request.Email,
		Name:    request.Name,
		IsAdmin: request.IsAdmin,
	}
	err := u.authenticationService.Register(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "Unable to register the user")
	}
	return &pb.EmptyResponse{}, nil
}

// UnregisterUser unregister a user
func (u userServer) UnregisterUser(ctx context.Context, request *pb.UserRequest) (*pb.EmptyResponse, error) {
	user := &entities.User{
		ID:      request.Id,
		Email:   request.Email,
		Name:    request.Name,
		IsAdmin: request.IsAdmin,
	}
	err := u.authenticationService.Unregister(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "Unable to unregister the user")
	}
	return &pb.EmptyResponse{}, nil
}

// Login and return user info and tokens
func (u userServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	response, err := u.authenticationService.Login(ctx, request.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Unable to login the user")
	}
	return &pb.LoginResponse{
		User: &pb.UserResponse{
			Id:      response.User.ID,
			Email:   response.User.Email,
			Name:    response.User.Name,
			IsAdmin: response.User.IsAdmin,
			Roles:   response.User.Roles,
		},
		Token: &pb.TokenResponse{
			Access:  response.Token.Access,
			Refresh: response.Token.Refresh,
		},
	}, nil
}

// Logout a user
func (u userServer) Logout(ctx context.Context, request *pb.LogoutRequest) (*pb.EmptyResponse, error) {
	err := u.authenticationService.Logout(ctx, &entities.Token{
		Access:  request.AccessToken,
		Refresh: request.RefreshToken,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Unable to logout")
	}
	return &pb.EmptyResponse{}, nil
}

// VerifyToken verify if a token is valid and return the associated user ID
func (u userServer) VerifyToken(ctx context.Context, request *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	ID, err := u.authenticationService.VerifyToken(ctx, request.AccessToken)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Invalid or expired token")
	}
	return &pb.VerifyResponse{
		Id: ID,
	}, nil
}

// RefreshToken the tokens
func (u userServer) RefreshToken(ctx context.Context, request *pb.RefreshRequest) (*pb.TokenResponse, error) {
	response, err := u.authenticationService.RefreshToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Invalid or expired refresh token")
	}
	return &pb.TokenResponse{
		Access:  response.Access,
		Refresh: response.Refresh,
	}, nil
}
