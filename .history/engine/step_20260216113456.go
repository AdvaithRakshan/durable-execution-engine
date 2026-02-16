package engine

import (
	"fmt"
)

func Step[T any](
	ctx *Context,
	stepID string,
	fn func() (T, error),
) (T, error) {

	var zero T

	// Generate unique step key using workflowID + stepID + sequence
	stepKey := ctx.NextStepKey(stepID)

	fmt.Println("[STEP] Checking:", stepKey)

	// Check if step already exists
	record, err := ctx.Persist.GetStep(stepKey)
	if err != nil {
		return zero, err
	}

	// If already completed → return cached result
	if record != nil && record.Status == "COMPLETED" {

		fmt.Println("[STEP] Skipping (cached):", stepKey)

		result, err := Deserialize[T](record.Result)
		if err != nil {
			return zero, err
		}

		return result, nil
	}

	// Insert PENDING record
	fmt.Println("[STEP] Executing:", stepKey)

	err = ctx.Persist.InsertPending(ctx.WorkflowID, stepKey)
	if err != nil {
		return zero, err
	}

	// Execute actual function
	result, err := fn()
	if err != nil {
		return zero, err
	}

	// Serialize result
	data, err := Serialize(result)
	if err != nil {
		return zero, err
	}

	// Mark step as completed
	err = ctx.Persist.MarkCompleted(stepKey, data)
	if err != nil {
		return zero, err
	}

	fmt.Println("[STEP] Completed:", stepKey)
	// Crash simulation
	if ctx.CrashAtStep == stepID {

		fmt.Println("[CRASH SIMULATION] Crashing at step:", stepID)

		panic("Simulated crash at " + stepID)
	}

	return result, nil
}
