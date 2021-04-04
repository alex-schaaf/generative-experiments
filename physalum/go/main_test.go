package main

import (
	"testing"
)

func TestMove(t *testing.T) {
	agent := Agent{100, 100, 90, 25, 9, 0, 0, 0}

	agent.move()

	if agent.y != 101 {
		t.Errorf("Agent.y is %f; want 101.", agent.y)
	}
}
