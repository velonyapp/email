//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"
	"log/slog"

	"github.com/velonyapp/notification/internal/application"
	"github.com/velonyapp/notification/internal/conf"
	"github.com/velonyapp/notification/internal/info"
	"github.com/velonyapp/notification/internal/infrastructure"
	"github.com/velonyapp/notification/internal/presentation"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(
	context.Context,
	*info.Service,
	*conf.Transport,
	*conf.Observability,
	*conf.Email,
	*slog.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		presentation.ProviderSet,
		infrastructure.ProviderSet,
		application.ProviderSet,
		newApp,
	))
}
