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
	// If step exists but is not completed, it means previous run crashed.
	// We must re-execute it safely.
	if record != nil && record.Status != "COMPLETED" {

		fmt.Println("[STEP] Found incomplete step, re-executing:", stepKey)

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
	// Crash simulation AFTER execution but BEFORE saving
	if ctx.CrashAtStep == stepID {

		fmt.Println("[CRASH SIMULATION] Crashing BEFORE commit:", stepID)

		panic("Simulated crash before commit at " + stepID)
	}

	fmt.Println("[STEP] Completed:", stepKey)
	// Crash simulation

	return result, nil
}
