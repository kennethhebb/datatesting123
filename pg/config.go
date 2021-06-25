package pg

import (
	"context"
	"fmt"
)

func (db Database) createConfigTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS config(version INTEGER NOT NULL DEFAULT 0);`
	if _, err := db.conn.Exec(ctx, query); err != nil {
		return fmt.Errorf("pg/Database.createConfigTable: %w", err)
	}
	return nil
}
