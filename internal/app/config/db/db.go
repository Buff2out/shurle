package db

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

func StartDB(driver string, DSN string) error {
	dbase, err := sql.Open(driver, DSN)
	if err != nil {
		return err
	}
	defer dbase.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = dbase.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
