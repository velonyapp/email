package mysql

import (
	"database/sql"
	"log/slog"

	"github.com/velonyapp/email/internal/conf"

	"github.com/XSAM/otelsql"
	"github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel/semconv/v1.37.0"
)

func NewConnection(c *conf.Data, log *slog.Logger) (*sql.DB, error) {
	cfg, err := mysql.ParseDSN(c.Mysql.Dsn)
	if err != nil {
		log.Error(
			"failed to parse MySQL DSN",
			slog.Any("error", err),
		)
		return nil, err
	}

	cfg.ParseTime = true
	cfg.InterpolateParams = true

	dsn := cfg.FormatDSN()

	attrs := append(otelsql.AttributesFromDSN(dsn), semconv.DBSystemNameMySQL)

	db, err := otelsql.Open("mysql", dsn,
		otelsql.WithAttributes(attrs...),
		otelsql.WithSpanOptions(otelsql.SpanOptions{
			OmitConnResetSession: true,
			OmitRows:             true,
		}),
	)
	if err != nil {
		log.Error(
			"failed to open MySQL connection",
			slog.Any("error", err),
		)
		return nil, err
	}

	if c.Mysql.MaxOpenConnections != nil {
		db.SetMaxOpenConns(int(*c.Mysql.MaxOpenConnections))
	}
	if c.Mysql.MaxIdleConnections != nil {
		db.SetMaxIdleConns(int(*c.Mysql.MaxIdleConnections))
	}
	if c.Mysql.MaxConnectionLifetime != nil {
		db.SetConnMaxLifetime(c.Mysql.MaxConnectionLifetime.AsDuration())
	}

	if err := db.Ping(); err != nil {
		log.Error(
			"failed to ping MySQL",
			slog.Any("error", err),
		)
		db.Close()
		return nil, err
	}

	log.Info(
		"MySQL connection established",
	)

	return db, nil
}
