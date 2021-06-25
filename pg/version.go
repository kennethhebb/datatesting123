package pg

import (
	"context"
	"fmt"
	"math"

	"github.com/jackc/pgx/v4"
)

func (db Database) GetCurrentVersion(ctx context.Context) (int, error) {
	if err := db.createConfigTable(ctx); err != nil {
		return 0, fmt.Errorf("pg/Database.GetCurrentVersion: %w", err)
	}

	version, err := db.getVersion(ctx)
	if err != nil {
		return 0, fmt.Errorf("pg/Database.GetCurrentVersion: %w", err)
	}

	return version, nil
}

func (db Database) ReconcileVersion(ctx context.Context, current, target int) error {
	// We take a 0 value for target to indicate that the database should be
	// migrated to the maximum version.
	if target == 0 {
		target = math.MaxInt64
	}

	var err error
	switch {
	case current < target:
		err = db.migrateUp(ctx, current, target)
	case target < current:
		err = db.migrateDown(ctx, current, target)
	default:
	}
	if err != nil {
		return fmt.Errorf("pg/Database.ReconcileVersion: %w", err)
	}
	return nil
}

func (db Database) getVersion(ctx context.Context) (int, error) {
	var version int
	err := db.conn.QueryRow(ctx, `SELECT version FROM config;`).Scan(&version)
	switch err {
	case nil:
		return version, nil
	case pgx.ErrNoRows:
		// We haven't initialized our config table yet, so we do so here. We'll
		// probably want to break this out into its own method once config gets
		// bigger.
		if _, err := db.conn.Exec(ctx, `INSERT INTO config (version) VALUES (0);`); err != nil {
			return 0, fmt.Errorf("pg/Database.getVersion: %w", err)
		}
		return 0, nil
	default:
		fmt.Println("is other err")
		return 0, fmt.Errorf("pg/Database.getVersion: %w", err)
	}
}

func (db Database) setVersion(ctx context.Context, version int) error {
	query := `UPDATE config SET version = $1;`
	if _, err := db.conn.Exec(ctx, query, version); err != nil {
		return fmt.Errorf("pg/Database.setVersion: %w", err)
	}
	return nil
}
