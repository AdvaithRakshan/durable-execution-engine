package main

import (
	"durable-engine/engine"
	"durable-engine/examples/onboarding"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("go run main.go <workflowID>")
		fmt.Println("go run main.go <workflowID> --crash-at <step>")
		return
	}

	workflowID := os.Args[1]

	crashAt := ""

	if len(os.Args) == 4 && os.Args[2] == "--crash-at" {
		crashAt = os.Args[3]
	}

	runner, err := engine.NewRunner("workflow.db")
	if err != nil {
		panic(err)
	}

	defer runner.Close()

	err = runner.Run(workflowID, crashAt, onboarding.OnboardingWorkflow)
	if err != nil {
		panic(err)
	}
}
