package pg

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func newMigration(path string) migration {
	parts := strings.Split(filepath.Base(path), "_")
	version, err := strconv.Atoi(parts[0])
	if err != nil {
		err = fmt.Errorf("pg/newMigration: migration path invalid: path=%s version=%s: %w", path, parts[0], err)
		panic(err)
	}
	return migration{
		Path:    path,
		Version: version,
	}
}

type migration struct {
	Path    string
	Version int
}

func (db Database) getMigrations(folder string) ([]migration, error) {
	path := filepath.Join(".", "migrations", folder, "*.sql")
	filenames, err := filepath.Glob(path)
	if err != nil {
		return nil, fmt.Errorf("pg/Database.getMigrations: %w", err)
	}

	migrations := make([]migration, 0, len(filenames))
	for _, fn := range filenames {
		migrations = append(migrations, newMigration(fn))
	}

	return migrations, nil
}

func (db Database) migrateDown(ctx context.Context, from, to int) error {
	// db.logger.WithField("from", from).WithField("to", to).Debug("migrating down")

	migrations, err := db.getMigrations("down")
	if err != nil {
		return fmt.Errorf("pg/Database.migrateDown: %w", err)
	}

	for i := len(migrations) - 1; i != 0; i-- {
		migration := migrations[i]
		// db.logger.WithField("version", migration.Version).Debug("checking migration")
		if migration.Version < to || migration.Version > from {
			continue
		}

		if err := db.execFile(ctx, migration.Path); err != nil {
			return fmt.Errorf("pg/Database.migrateDown: %w", err)
		}
	}

	if err := db.setVersion(ctx, to); err != nil {
		return fmt.Errorf("pg/Database.migrateDown: %w", err)
	}

	return nil
}

func (db Database) migrateUp(ctx context.Context, from, to int) error {
	// db.logger.WithField("from", from).WithField("to", to).Debug("migrating up")

	migrations, err := db.getMigrations("up")
	if err != nil {
		return fmt.Errorf("pg/Database.migrateUp: %w", err)
	}

	// "To" can be greater than the max possible version, such as when the user
	// has not specified a particular target in which case "to" will equal the
	// maximum possible int64. We need to clamp that to our actual max version.
	if to > len(migrations) {
		to = len(migrations)
	}

	// TODO: The file migrations and the version setting should all be run in
	// a transaction
	for _, migration := range migrations {
		// Outside the range of migrations to run. We don't re-run previous
		// migrations since they may duplicate data or cause errors where they
		// shouldn't.
		if migration.Version > to || migration.Version <= from {
			continue
		}

		if err := db.execFile(ctx, migration.Path); err != nil {
			return fmt.Errorf("pg/Database.migrateUp: %w", err)
		}
	}

	if err := db.setVersion(ctx, to); err != nil {
		return fmt.Errorf("pg/Database.migrateUp: %w", err)
	}

	return nil
}
