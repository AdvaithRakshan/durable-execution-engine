package engine

import (
	"fmt"
	"sync/atomic"
)

type SequenceManager struct {
	counter atomic.Int64
}

func NewSequenceManager() *SequenceManager {
	return &SequenceManager{}
}

func (s *SequenceManager) Next(stepID string) string {

	seq := s.counter.Add(1)

	return fmt.Sprintf("%s-%d", stepID, seq)
}
