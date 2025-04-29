//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"learn/check_status/client"
	"learn/check_status/config"
	"learn/check_status/controller"
	"learn/check_status/dao"
	"learn/check_status/route"
	"learn/check_status/service"
	"learn/check_status/tool"
)

func InitApp(cfg string) (*route.App, error) {
	wire.Build(
		route.ProviderSet,
		controller.ProviderSet,
		service.ProviderSet,
		client.ProviderSet,
		config.ProviderSet,
		dao.ProviderSet,
		tool.ProviderSet,
		wire.Bind(new(service.CheckService), new(*service.Check)),
		wire.Bind(new(dao.CheckDAO), new(*dao.CheckDAOImpl)),
	)
	return nil, nil
}
