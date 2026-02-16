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
func (p *Persistence) GetStep(stepKey string) (*StepRecord, error) {

	row := p.DB.QueryRow(
		`SELECT workflow_id, step_key, status, result 
		 FROM steps WHERE step_key = ?`,
		stepKey,
	)

	var record StepRecord

	err := row.Scan(
		&record.WorkflowID,
		&record.StepKey,
		&record.Status,
		&record.Result,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &record, nil
}
func (p *Persistence) InsertPending(
	workflowID string,
	stepKey string,
) error {

	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	_, err := p.DB.Exec(
		`INSERT INTO steps(workflow_id, step_key, status)
         VALUES(?, ?, ?)
         ON CONFLICT(step_key) DO UPDATE SET status='PENDING', result=NULL`,
		workflowID,
		stepKey,
		"PENDING",
	)

	return err
}

func (p *Persistence) MarkCompleted(
	stepKey string,
	result []byte,
) error {

	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	_, err := p.DB.Exec(
		`UPDATE steps
		 SET status = ?, result = ?
		 WHERE step_key = ?`,
		"COMPLETED",
		result,
		stepKey,
	)

	return err
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
