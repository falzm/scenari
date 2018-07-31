// Package scenari provides a simple interface to describe and execute arbitrary scenarios.
package scenari

import (
	"container/list"
	"time"
)

// Scenario represents a scenario instance.
type Scenario struct {
	name     string
	steps    *list.List
	carryOn  bool
	rep      int
	repDelay time.Duration
}

// NewScenario returns an new scenario instance.
func NewScenario(name string) *Scenario {
	return &Scenario{
		name:  name,
		steps: list.New(),
		rep:   1,
	}
}

// Step adds a new step in the scenario.
func (s *Scenario) Step(step *Step) *Scenario {
	s.steps.PushBack(step)

	return s
}

// Pause injects a pause of duration d after the latest step added to the scenario.
func (s *Scenario) Pause(d time.Duration) *Scenario {
	return s.Step(NewStep(func() error {
		time.Sleep(d)
		return nil
	}))
}

// CarryOn sets the scenario rollout to continue even if one ore more steps return an error.
func (s *Scenario) CarryOn() *Scenario {
	s.carryOn = true

	return s
}

// Repeat sets the scenario to be repeated n times with a delay of duration d between each iteration.
func (s *Scenario) Repeat(n int, delay time.Duration) *Scenario {
	s.rep = n
	s.repDelay = delay

	return s
}

// Rollout executes the scenario.
func (s *Scenario) Rollout() error {
	var err error

	for i := 0; i < s.rep; i++ {
		if i > 0 && s.repDelay > 0 {
			time.Sleep(s.repDelay)
		}

		for step := s.steps.Front(); step != nil; step = step.Next() {
			if err = step.Value.(*Step).Exec(); err != nil && !s.carryOn {
				return err
			}
		}
	}

	return nil
}
