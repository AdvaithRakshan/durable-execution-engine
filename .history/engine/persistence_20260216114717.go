package engine

import (
	"database/sql"
	"encoding/json"
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

	_, err = db.Exec(`PRAGMA journal_mode=WAL;`)
	if err != nil {
		return nil, err
	}

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

	return &Persistence{DB: db}, nil
}

func Serialize[T any](value T) ([]byte, error) {
	return json.Marshal(value)
}

func Deserialize[T any](data []byte) (T, error) {

	var value T

	err := json.Unmarshal(data, &value)

	return value, err
}

func (p *Persistence) Close() error {
	return p.DB.Close()
}
