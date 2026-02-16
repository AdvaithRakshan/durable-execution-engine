package main

import (
	"durable-engine/engine"
	"fmt"
)

func main() {

	persist, err := engine.NewPersistence("workflow.db")
	if err != nil {
		panic(err)
	}

	ctx := engine.NewContext("workflow1", persist)

	result, err := engine.Step(ctx, "test-step", func() (string, error) {
		fmt.Println("Executing actual function")
		return "hello world", nil
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Result:", result)

	defer persist.Close()
}
