package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateResourceTable, downCreateResourceTable)
}

func upCreateResourceTable(tx *sql.Tx) error {
	_, err := tx.Exec(
		"CREATE TABLE IF NOT EXISTS resource( " +
			"id BIGSERIAL NOT NULL PRIMARY KEY," +
			"user_id BIGINT NOT NULL," +
			"type BIGINT NOT NULL," +
			"status BIGINT NOT NULL" +
			");",
	)
	return err
}

func downCreateResourceTable(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE resource;")
	return err
}
