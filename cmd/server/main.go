package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/velonyapp/email/internal/conf"
	"github.com/velonyapp/email/internal/info"
	"github.com/velonyapp/email/internal/infrastructure/observability"
	"github.com/velonyapp/email/internal/infrastructure/transport"

	"buf.build/go/protovalidate"
	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
	"github.com/go-kratos/kratos/v3"
	"github.com/go-kratos/kratos/v3/config"
	"github.com/go-kratos/kratos/v3/config/env"
	"github.com/go-kratos/kratos/v3/config/file"
	"github.com/go-kratos/kratos/v3/log"
	"github.com/go-kratos/kratos/v3/transport/grpc"
	"github.com/go-kratos/kratos/v3/transport/http"

	_ "go.uber.org/automaxprocs"
)

var (
	Name          = "velony-email"
	Version       = "dev"
	InstanceID, _ = os.Hostname()

	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "config", "../../configs", "config path, eg: -config config.yaml")
}

func newApp(
	logger *slog.Logger,
	gs *grpc.Server,
	hs *http.Server,
	rs *transport.RabbitMQServer,
	_ *observability.OpenTelemetry,
) *kratos.App {
	return kratos.New(
		kratos.ID(InstanceID),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
			rs,
		),
	)
}

func main() {
	flag.Parse()

	// Info
	bi := info.Bootstrap{
		Service: &info.Service{
			Name:       Name,
			Version:    Version,
			InstanceID: InstanceID,
		},
	}
	if err := protovalidate.Validate(&bi); err != nil {
		panic(err)
	}

	// Config
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
			env.NewSource("VELONY_EMAIL_"),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	if err := protovalidate.Validate(&bc); err != nil {
		panic(err)
	}

	// Logger
	logger := log.NewLogger(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}),
		log.WithExtractor(tracing.TraceAttrs),
	).With(
		slog.String("service.name", Name),
		slog.String("service.version", Version),
		slog.String("service.instance.id", InstanceID),
	)
	log.SetDefault(logger)

	// App
	app, cleanup, err := wireApp(
		context.Background(),
		bi.Service,
		bc.Email,
		bc.Transport,
		bc.Observability,
		logger,
	)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
