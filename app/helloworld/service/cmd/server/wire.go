//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/ray-dota/backend-mono/app/helloworld/service/internal/biz"
	"github.com/ray-dota/backend-mono/app/helloworld/service/internal/conf"
	"github.com/ray-dota/backend-mono/app/helloworld/service/internal/data"
	"github.com/ray-dota/backend-mono/app/helloworld/service/internal/server"
	"github.com/ray-dota/backend-mono/app/helloworld/service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
