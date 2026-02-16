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
