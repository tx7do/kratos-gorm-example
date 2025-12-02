//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"kratos-gorm-example/api/gen/go/common/conf"

	"github.com/google/wire"

	"kratos-gorm-example/app/user/service/internal/data"
	"kratos-gorm-example/app/user/service/internal/server"
	"kratos-gorm-example/app/user/service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

// initApp init kratos application.
func initApp(log.Logger, *conf.Bootstrap) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, service.ProviderSet, newApp))
}
