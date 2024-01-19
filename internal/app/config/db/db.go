package db

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

func StartDB(driver string, DSN string) (*sql.DB, error) {
	dbase, err := sql.Open(driver, DSN)
	if err != nil {
		return nil, err
	}
	defer dbase.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = dbase.PingContext(ctx); err != nil {
		return nil, err
	}
	return dbase, nil
}
