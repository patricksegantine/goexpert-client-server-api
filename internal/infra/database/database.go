package database

import "database/sql"

const createTable string = `
	CREATE TABLE IF NOT EXISTS exchange_rate (
		ID 			INTEGER NOT NULL PRIMARY KEY,
		name       	VARCHAR,
		bid        	VARCHAR,
		created_at 	VARCHAR
	);
`

func SetupDatabase(db *sql.DB) error {
	if _, err := db.Exec(createTable); err != nil {
		return err
	}
	return nil
}
