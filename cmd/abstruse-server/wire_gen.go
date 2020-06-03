// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jkuri/abstruse/internal/pkg/auth"
	"github.com/jkuri/abstruse/internal/pkg/config"
	"github.com/jkuri/abstruse/internal/pkg/http"
	"github.com/jkuri/abstruse/internal/pkg/log"
	"github.com/jkuri/abstruse/internal/server"
	"github.com/jkuri/abstruse/internal/server/app"
	"github.com/jkuri/abstruse/internal/server/controller"
	"github.com/jkuri/abstruse/internal/server/db"
	"github.com/jkuri/abstruse/internal/server/db/repository"
	"github.com/jkuri/abstruse/internal/server/etcd"
	"github.com/jkuri/abstruse/internal/server/service"
	"github.com/jkuri/abstruse/internal/server/websocket"
)

// Injectors from wire.go:

func CreateApp(cfg string) (*server.App, error) {
	viper, err := config.NewConfig(cfg)
	if err != nil {
		return nil, err
	}
	options, err := server.NewOptions(viper)
	if err != nil {
		return nil, err
	}
	logOptions, err := log.NewOptions(viper)
	if err != nil {
		return nil, err
	}
	logger, err := log.New(logOptions)
	if err != nil {
		return nil, err
	}
	httpOptions, err := http.NewOptions(viper)
	if err != nil {
		return nil, err
	}
	websocketOptions, err := websocket.NewOptions(viper)
	if err != nil {
		return nil, err
	}
	dbOptions, err := db.NewOptions(viper, logger)
	if err != nil {
		return nil, err
	}
	gormDB, err := db.NewDatabase(dbOptions)
	if err != nil {
		return nil, err
	}
	userRepository := repository.NewDBUserRepository(logger, gormDB)
	userService := service.NewUserService(logger, userRepository)
	userController := controller.NewUserController(logger, userService)
	versionService := service.NewVersionService(logger)
	versionController := controller.NewVersionController(logger, versionService)
	appOptions, err := app.NewOptions(viper)
	if err != nil {
		return nil, err
	}
	websocketApp := websocket.NewApp(logger)
	repoRepository := repository.NewDBRepoRepository(gormDB)
	jobRepository := repository.NewDBJobRepository(gormDB)
	buildRepository := repository.NewDBBuildRepository(gormDB)
	appApp, err := app.NewApp(appOptions, websocketApp, repoRepository, jobRepository, buildRepository, logger)
	if err != nil {
		return nil, err
	}
	workerService := service.NewWorkerService(logger, appApp)
	workerController := controller.NewWorkerController(logger, workerService)
	buildService := service.NewBuildService(buildRepository, jobRepository, appApp)
	buildController := controller.NewBuildController(buildService)
	repositoryService := service.NewRepositoryService(repoRepository)
	repositoryController := controller.NewRepositoryController(repositoryService)
	providerRepository := repository.NewDBProviderRepository(gormDB)
	providerService := service.NewProviderService(providerRepository)
	providerController := controller.NewProviderController(providerService, repositoryService)
	middlewareController := controller.NewMiddlewareController(logger, userService)
	initControllers := controller.CreateInitControllersFn(userController, versionController, workerController, buildController, repositoryController, providerController, middlewareController)
	router := http.NewRouter(httpOptions, websocketOptions, initControllers)
	httpServer, err := http.NewServer(httpOptions, logger, router)
	if err != nil {
		return nil, err
	}
	etcdOptions, err := etcd.NewOptions(viper)
	if err != nil {
		return nil, err
	}
	etcdServer := etcd.NewServer(etcdOptions, logger)
	websocketServer := websocket.NewServer(websocketOptions, websocketApp, logger)
	serverApp := server.NewApp(options, logger, httpServer, etcdServer, websocketServer, appApp)
	return serverApp, nil
}

// wire.go:

var providerSet = wire.NewSet(log.ProviderSet, config.ProviderSet, http.ProviderSet, etcd.ProviderSet, app.ProviderSet, db.ProviderSet, repository.ProviderSet, auth.ProviderSet, controller.ProviderSet, service.ProviderSet, websocket.ProviderSet, server.ProviderSet)
