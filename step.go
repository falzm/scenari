package scenari

// Step reprensents a scenario step instance.
type Step struct {
	preExecFunc  func() error
	execFunc     func() error
	postExecFunc func() error
}

// NewStep returns a new scenario step instance.
func NewStep(f func() error) *Step {
	return &Step{execFunc: f}
}

// PreExec defines a function f to be executed before the step function.
func (s *Step) PreExec(f func() error) *Step {
	s.preExecFunc = f

	return s
}

// PostExec defines a function f to be executed after the step function.
func (s *Step) PostExec(f func() error) *Step {
	s.postExecFunc = f

	return s
}

// Exec executes the step function, and its pre/post functions if any.
func (s *Step) Exec() error {
	var err error

	if s.preExecFunc != nil {
		if err = s.preExecFunc(); err != nil {
			return err
		}
	}

	if err = s.execFunc(); err != nil {
		return err
	}

	if s.postExecFunc != nil {
		if err = s.postExecFunc(); err != nil {
			return err
		}
	}

	return nil
}
