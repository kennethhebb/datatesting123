package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kennethhebb/datatesting123/pg"
)

func main() {
	var cfg Config
	if err := cfg.Load(os.Args[1:]); err != nil {
		panic(err)
	}
	if err := Run(cfg); err != nil {
		panic(err)
	}
}

func Run(cfg Config) error {
	ctx := context.Background()

	db, err := pg.NewDatabase(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("Run.NewDatabase: %w", err)
	}

	currentVersion, err := db.GetCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}

	if err := db.ReconcileVersion(ctx, currentVersion, cfg.DatabaseVersion); err != nil {
		return fmt.Errorf("Run: %w", err)
	}

	return nil

}
