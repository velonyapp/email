//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	"github.com/velonyapp/email/internal/application"
	"github.com/velonyapp/email/internal/conf"
	"github.com/velonyapp/email/internal/domain"
	"github.com/velonyapp/email/internal/infrastructure"
	"github.com/velonyapp/email/internal/presentation"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(
	*conf.Email,
	*conf.Data,
	*conf.Transport,
	*slog.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		domain.ProviderSet,
		application.ProviderSet,
		infrastructure.ProviderSet,
		presentation.ProviderSet,
		newApp,
	))
}
