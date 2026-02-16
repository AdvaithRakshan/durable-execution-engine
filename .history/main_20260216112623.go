package main

import (
	"durable-engine/engine"
	"durable-engine/examples/onboarding"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <workflowID>")
		return
	}

	workflowID := os.Args[1]

	runner, err := engine.NewRunner("workflow.db")
	if err != nil {
		panic(err)
	}

	defer runner.Close()

	err = runner.Run(workflowID, onboarding.OnboardingWorkflow)
	if err != nil {
		panic(err)
	}
}
