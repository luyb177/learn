//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"learn/Grab_seat/client"
	"learn/Grab_seat/config"
	"learn/Grab_seat/controller"
	"learn/Grab_seat/dao"
	"learn/Grab_seat/route"
	"learn/Grab_seat/service"
)

func InitApp(ConfigPath string) (*route.App, error) {
	wire.Build(
		route.ProviderSet,
		controller.ProviderSet,
		service.ProviderSet,
		dao.ProviderSet,
		config.ProviderSet,
		client.ProviderSet,
		wire.Bind(new(service.GrabService), new(*service.GrabServiceImpl)),
		wire.Bind(new(service.MonitorService), new(*service.MonitorServiceImpl)),
		wire.Bind(new(service.SseService), new(*service.SseServiceImpl)),
		wire.Bind(new(dao.GrabDAO), new(*dao.GrabDAOImpl)),
		wire.Bind(new(dao.ContentDAO), new(*dao.ContentDAOImpl)),
	)

	return &route.App{}, nil
}
