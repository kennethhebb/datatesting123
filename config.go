package main

import (
	"github.com/kennethhebb/datatesting123/lib/flag"

	"github.com/kennethhebb/datatesting123/lib/os"
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
		os.GetIntEnv("DATABASE_VERSION"),
		"The database version to use",
	)
	fs.Parse(args)

	return nil
}
