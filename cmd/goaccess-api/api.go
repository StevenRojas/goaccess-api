package main

import (
	"context"
	"fmt"
	"github.com/StevenRojas/goaccess-api/pkg/pb"
	"github.com/StevenRojas/goaccess-api/pkg/server"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/StevenRojas/goaccess-api/pkg/transport"
	"github.com/StevenRojas/goaccess/pkg/configuration"
	"github.com/StevenRojas/goaccess/pkg/service"
	"github.com/gorilla/mux"
	"github.com/oklog/oklog/pkg/group"
)

var (
	authenticationService service.AuthenticationService
	accessService         service.AccessService
	authorizationService  service.AuthorizationService
	initService           service.InitializationService
)

func main() {
	ctx := context.Background()
	serviceConfig, err := configuration.Read()
	if err != nil {
		panic(err)
	}

	logger := configuration.NewLogger(serviceConfig.Server)
	logger.Debug("creating services...")
	factory := service.NewServiceFactory(ctx, serviceConfig)
	factory.Setup()
	authenticationService = factory.CreateAuthenticationService()
	accessService = factory.CreateAccessService()
	authorizationService = factory.CreateAuthorizationService()
	initService = factory.CreateInitializationService()
	logger.Debug("services ready")

	router := mux.NewRouter()
	transport.MakeHTTPHandlerForAccess(router, accessService, serviceConfig.Security, logger)
	transport.MakeHTTPHandlerForActions(router, authorizationService, serviceConfig.Security, logger)
	transport.MakeHTTPHandlerForInit(router, initService, serviceConfig.Security, logger)

	go setGRPC(serviceConfig.Server.GRPC)
	logger.Info("GRPC server listen at " + serviceConfig.Server.GRPC)

	var runGroup group.Group
	{
		httpServer := http.Server{
			Addr:    serviceConfig.Server.HTTP,
			Handler: router,
		}
		runGroup.Add(func() error {
			logger.Info("HTTP server listen at " + serviceConfig.Server.HTTP)
			return httpServer.ListenAndServe() // TODO: support TLS
		}, func(err error) {
			httpServer.Shutdown(context.Background())
			logger.Error("HTTP server shutdown with error", "error", err)
		})
	}

	{
		cancel := make(chan struct{})
		runGroup.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case s := <-c:
				return fmt.Errorf("signal received %s", s)
			case <-cancel:
				return nil
			}
		}, func(error) {
			close(cancel)
		})
	}
	runGroup.Run()
	logger.Info("server terminated")
}

func setGRPC(port string) {
	grpcServer := grpc.NewServer()
	pb.RegisterUsersServer(grpcServer, server.NewUserServer(authenticationService))
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
