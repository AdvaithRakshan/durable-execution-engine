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

	fmt.Println(ctx.NextStepKey("step"))
	fmt.Println(ctx.NextStepKey("step"))

	defer persist.Close()
}
