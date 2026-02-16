package engine

import (
	"sync"
)

type Context struct {
	WorkflowID string
	Persist    *Persistence
	Sequence   *SequenceManager
	Mutex      sync.Mutex

	CrashAtStep string
}

func NewContext(
	workflowID string,
	persist *Persistence,
) *Context {

	return &Context{
		WorkflowID: workflowID,
		Persist:    persist,
		Sequence:   NewSequenceManager(),
	}
}
func (c *Context) NextStepKey(stepID string) string {

	return c.Sequence.Next(c.WorkflowID, stepID)

}
