package engine

import (
	"database/sql"
	"fmt"
)

func Step[T any](
	ctx *Context,
	stepID string,
	fn func() (T, error),
) (T, error) {

	var zero T

	stepKey := ctx.NextStepKey(stepID)

	fmt.Println("[STEP] Checking:", stepKey)

	// Lock persistence for transaction safety
	ctx.Persist.Mutex.Lock()

	tx, err := ctx.Persist.DB.Begin()
	if err != nil {
		ctx.Persist.Mutex.Unlock()
		return zero, err
	}

	// CRITICAL: Always release resources, even on panic
	defer func() {
		tx.Rollback()
		ctx.Persist.Mutex.Unlock()
	}()

	var status string
	var resultData []byte

	err = tx.QueryRow(
		`SELECT status, result FROM steps WHERE step_key = ?`,
		stepKey,
	).Scan(&status, &resultData)

	// CASE 1: Completed → return cached
	if err == nil && status == "COMPLETED" {

		fmt.Println("[STEP] Skipping (cached):", stepKey)

		result, err := Deserialize[T](resultData)
		if err != nil {
			return zero, err
		}

		return result, nil
	}

	// CASE 2: New step → insert pending
	if err == sql.ErrNoRows {

		_, err = tx.Exec(
			`INSERT INTO steps(workflow_id, step_key, status)
			 VALUES(?, ?, ?)`,
			ctx.WorkflowID,
			stepKey,
			"PENDING",
		)

		if err != nil {
			return zero, err
		}

	} else if err == nil && status == "PENDING" {

		fmt.Println("[STEP] Recovering incomplete step:", stepKey)

	} else if err != nil && err != sql.ErrNoRows {
		return zero, err
	}

	fmt.Println("[STEP] Executing:", stepKey)

	result, err := fn()
	if err != nil {
		return zero, err
	}

	// Crash simulation BEFORE commit
	if ctx.CrashAtStep == stepID {

		fmt.Println("[CRASH SIMULATION] Crashing BEFORE commit:", stepID)

		panic("Simulated crash before commit at " + stepID)
	}

	data, err := Serialize(result)
	if err != nil {
		return zero, err
	}

	_, err = tx.Exec(
		`UPDATE steps
		 SET status = ?, result = ?
		 WHERE step_key = ?`,
		"COMPLETED",
		data,
		stepKey,
	)

	if err != nil {
		return zero, err
	}

	err = tx.Commit()
	if err != nil {
		return zero, err
	}

	fmt.Println("[STEP] Completed:", stepKey)

	return result, nil
}
