package main

import (
	"durable-engine/engine"
	"fmt"
)

func main() {

	seq := engine.NewSequenceManager()

	fmt.Println(seq.Next("workflow1", "step"))
	fmt.Println(seq.Next("workflow1", "step"))
	fmt.Println(seq.Next("workflow1", "step"))
}
