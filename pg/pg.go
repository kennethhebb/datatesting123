package pg

import (
	"context"
	"fmt"
	"os"

	pgx "github.com/jackc/pgx/v4"
)

type Option func(*Database)

type Tx struct {
	tx pgx.Tx
}

func NewDatabase(ctx context.Context, url string, opts ...Option) (*Database, error) {
	fmt.Println("this! " + url)
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("pg/NewDatabase: %w", err)
	}
	fmt.Println(conn.Config().ConnString())
	db := &Database{
		conn: conn,
		//logger: log.NewLogger(), // TODO: This should log nothing if not provided
	}
	for _, opt := range opts {
		opt(db)
	}
	return db, nil
}

type Database struct {
	conn *pgx.Conn
	//logger log.Logger
}

func (db Database) execFile(ctx context.Context, path string) error {
	//db.logger.WithField("path", path).Debug("running migration")

	migration, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("pg/Database.ExecFile: %w", err)
	}

	if _, err := db.conn.Exec(ctx, string(migration)); err != nil {
		return fmt.Errorf("pg/Database.ExecFile: %w", err)
	}

	return nil
}

func (db *Database) Close(ctx context.Context) error {
	if err := db.conn.Close(ctx); err != nil {
		return fmt.Errorf("pg/Database.Close: %w", err)
	}
	return nil
}
