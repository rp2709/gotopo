package gotopo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	pgxgeom "github.com/twpayne/pgx-geom"
)

type DatabaseConfig struct {
	Dbuser string
	Dbname string
	Dbpwd  string
	Dbhost string
	Dbport int
}

const connectionStringFmt = "user=%s dbname=%s sslmode=prefer host=%s password=%s port=%d"

func (config DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		connectionStringFmt,
		config.Dbuser,
		config.Dbname,
		config.Dbhost,
		config.Dbpwd,
		config.Dbport,
	)
}

func (config DatabaseConfig) Connect(ctx context.Context) (*Queries, *pgx.Conn, error) {
	// do not add the database password or username to log
	slog.Info(
		"Connecting to database",
		"host", config.Dbhost,
		"port", config.Dbport,
		"database", config.Dbname,
	)

	conn, err := pgx.Connect(ctx, config.ConnectionString())
	if err != nil {
		return nil, nil, err
	}

	err = pgxgeom.Register(ctx, conn)
	if err != nil {
		return nil, nil, err
	}

	queries := New(conn)
	return queries, conn, nil
}
