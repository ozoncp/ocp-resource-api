package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/ozoncp/ocp-resource-api/cmd/migration/migrations"
	"github.com/pressly/goose/v3"
)

const directory = "."

var (
	flags = flag.NewFlagSet("migration", flag.ExitOnError)
)

func main() {
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatal("migration: failed to parse args\n", err)
	}
	args := flags.Args()

	if len(args) < 2 {
		flags.Usage()
		return
	}

	dbstring, command := args[0], args[1]
	db, err := sql.Open("pgx", dbstring)
	if err != nil {
		log.Fatalf("migration failed to open: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("migration failed to close: %v\n", err)
		}
	}()

	arguments := make([]string, 0, len(args))
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	arguments = []string{}

	if err := goose.Run(command, db, directory, arguments...); err != nil {
		log.Fatalf("migration %v: %v", command, err)
	}
}
