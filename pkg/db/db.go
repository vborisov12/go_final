package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX idx_date ON scheduler(date);
`

type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(db *sql.DB) TaskStore {
	return TaskStore{db: db}
}

func Init(dbPath string) error {
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		log.Printf("Creating database at %s", dbPath)
		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			return err
		}
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
		db.Close()
	} else {
		log.Printf("Database already exists at %s", dbPath)
	}
	return nil
}
