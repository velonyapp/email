package mysql

import (
	"database/sql"

	"github.com/velonyapp/email/internal/conf"

	"github.com/XSAM/otelsql"
	"github.com/go-sql-driver/mysql"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func NewConnection(c *conf.Data) (*sql.DB, error) {
	cfg, err := mysql.ParseDSN(c.Mysql.Dsn)
	if err != nil {
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
		db.Close()
		return nil, err
	}

	return db, nil
}
