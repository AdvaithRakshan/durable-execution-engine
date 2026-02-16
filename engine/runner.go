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

func (r *Runner) Run(
	workflowID string,
	crashAt string,
	workflow WorkflowFunc,
) error {

	fmt.Println("[RUNNER] Starting workflow:", workflowID)

	if crashAt != "" {
		fmt.Println("[RUNNER] Crash simulation enabled at step:", crashAt)
	}

	ctx := NewContext(workflowID, r.Persist, crashAt)

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
