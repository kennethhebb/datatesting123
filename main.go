package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kennethhebb/datatesting123/lib/flag"
	osx "github.com/kennethhebb/datatesting123/lib/os"
	"github.com/kennethhebb/datatesting123/pg"
)

type Config struct {
	DatabaseURL     string
	DatabaseVersion int
}

func (cfg *Config) Load(args []string) error {
	var printHelp bool

	fs := flag.NewFlagSet("server")
	fs.BoolVar(
		&printHelp,
		"h",
		"help",
		false,
		"Print help information",
	)
	fs.StringVar(
		&cfg.DatabaseURL,
		"db",
		"database-url",
		//os.GetStringEnv("DATABASE_URL"),
		"postgres://postgres:123456789@localhost:5434/cloud?sslmode=disable",
		"The URL of a database to connect to",
	)
	fs.IntVar(
		&cfg.DatabaseVersion,
		"",
		"database-version",
		osx.GetIntEnv("DATABASE_VERSION"),
		"The database version to use",
	)
	fs.Parse(args)

	return nil
}

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
