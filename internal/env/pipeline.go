package env

import "fmt"

// StepFunc is a transformation applied to a Set during a pipeline run.
type StepFunc func(*Set) (*Set, error)

// Pipeline represents an ordered sequence of transformation steps applied to a Set.
type Pipeline struct {
	name  string
	steps []StepFunc
}

// NewPipeline creates a new Pipeline with the given name.
func NewPipeline(name string) (*Pipeline, error) {
	if name == "" {
		return nil, fmt.Errorf("pipeline name must not be empty")
	}
	return &Pipeline{name: name}, nil
}

// Name returns the pipeline name.
func (p *Pipeline) Name() string { return p.name }

// AddStep appends a StepFunc to the pipeline.
func (p *Pipeline) AddStep(fn StepFunc) error {
	if fn == nil {
		return fmt.Errorf("step func must not be nil")
	}
	p.steps = append(p.steps, fn)
	return nil
}

// Len returns the number of steps in the pipeline.
func (p *Pipeline) Len() int { return len(p.steps) }

// Run executes each step in order, passing the output of one step as the input
// to the next. The original set is not modified; a copy is made before the
// first step runs.
func (p *Pipeline) Run(src *Set) (*Set, error) {
	if src == nil {
		return nil, fmt.Errorf("source set must not be nil")
	}
	current, err := Clone(src, src.Name())
	if err != nil {
		return nil, fmt.Errorf("pipeline clone: %w", err)
	}
	for i, step := range p.steps {
		current, err = step(current)
		if err != nil {
			return nil, fmt.Errorf("pipeline step %d: %w", i, err)
		}
		if current == nil {
			return nil, fmt.Errorf("pipeline step %d returned nil set", i)
		}
	}
	return current, nil
}
