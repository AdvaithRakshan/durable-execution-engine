package engine

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Persistence struct {
	DB    *sql.DB
	Mutex sync.Mutex
}

type StepRecord struct {
	WorkflowID string
	StepKey    string
	Status     string
	Result     []byte
}

func NewPersistence(dbPath string) (*Persistence, error) {

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Enable WAL mode for concurrency
	_, err = db.Exec(`PRAGMA journal_mode=WAL;`)
	if err != nil {
		return nil, err
	}

	// Prevent SQLITE_BUSY errors
	_, err = db.Exec(`PRAGMA busy_timeout = 5000;`)
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS steps (
		workflow_id TEXT,
		step_key TEXT PRIMARY KEY,
		status TEXT,
		result BLOB,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &Persistence{
		DB: db,
	}, nil
}
