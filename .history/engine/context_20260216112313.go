package engine

import (
	"sync"
)

type Context struct {
	WorkflowID string
	Persist    *Persistence
	Sequence   *SequenceManager
	Mutex      sync.Mutex
}
