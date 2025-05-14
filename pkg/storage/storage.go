package storage

import (
	"database/sql"
	"fmt"
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

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskStore struct {
	db *sql.DB
}

func NewTaskStorage(db *sql.DB) TaskStore {
	return TaskStore{db: db}
}

func Init(dbPath string) (*sql.DB, error) {
	needInit := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("Create new database %s", dbPath)
		needInit = true
	} else if err != nil {
		return nil, fmt.Errorf("Database check failed: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("Database open failed: %w", err)
	}

	if needInit {
		if _, err := db.Exec(schema); err != nil {
			return nil, fmt.Errorf("Database init failed: %w", err)
		}
		log.Println("Database initialized")
	} else {
		log.Printf("Database already exists: %s", dbPath)
	}

	return db, nil
}
