package main

import (
	"durable-engine/engine"
	"fmt"
)

func main() {

	p, err := engine.NewPersistence("workflow.db")
	if err != nil {
		panic(err)
	}

	fmt.Println("Database initialized")

	defer p.Close()
}
