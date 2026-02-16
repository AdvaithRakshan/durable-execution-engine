package engine

import (
	"fmt"
)

type WorkflowFunc func(ctx *Context) error

type Runner struct {
	Persist *Persistence
}

func NewRunner(dbPath string) (*Runner, error) {

	persist, err := NewPersistence(dbPath)
	if err != nil {
		return nil, err
	}

	return &Runner{
		Persist: persist,
	}, nil
}

func (r *Runner) Run(workflowID string, workflow WorkflowFunc) error {

	fmt.Println("[RUNNER] Starting workflow:", workflowID)

	ctx := NewContext(workflowID, r.Persist)

	err := workflow(ctx)
	if err != nil {
		fmt.Println("[RUNNER] Workflow failed:", err)
		return err
	}

	fmt.Println("[RUNNER] Workflow completed:", workflowID)

	return nil
}

func (r *Runner) Close() error {
	return r.Persist.Close()
}
