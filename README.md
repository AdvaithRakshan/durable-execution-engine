# Durable Execution Engine (Go + SQLite)

## Overview

This project implements a **Durable Execution Engine** in Go that allows workflows written in normal native Go code to survive crashes and resume execution from the exact point of failure.

Unlike traditional programs where a crash resets execution, this engine ensures that completed steps are never re-executed and incomplete steps are safely resumed.

This design is inspired by durable execution systems such as:

- Temporal
- Cadence
- Azure Durable Functions
- DBOS

---

## Key Features

### Durable Execution
Each step result is persisted in SQLite. If the workflow crashes, completed steps are skipped and execution resumes from the last incomplete step.

### Generic Step Primitive
Supports any return type using Go generics:

```go
func Step[T any](ctx *Context, stepID string, fn func() (T, error)) (T, error)
