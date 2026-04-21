package env

import (
	"fmt"
	"strings"
	"testing"
)

func makePipelineSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("base")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	_ = s.Put("HOST", "localhost")
	_ = s.Put("PORT", "8080")
	_ = s.Put("debug", "true")
	return s
}

func TestNewPipelineEmptyNameReturnsError(t *testing.T) {
	_, err := NewPipeline("")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNewPipelineCreatesEmptySteps(t *testing.T) {
	p, err := NewPipeline("mypipe")
	if err != nil {
		t.Fatalf("NewPipeline: %v", err)
	}
	if p.Len() != 0 {
		t.Fatalf("expected 0 steps, got %d", p.Len())
	}
	if p.Name() != "mypipe" {
		t.Fatalf("expected name 'mypipe', got %q", p.Name())
	}
}

func TestAddStepNilFuncReturnsError(t *testing.T) {
	p, _ := NewPipeline("p")
	if err := p.AddStep(nil); err == nil {
		t.Fatal("expected error for nil step")
	}
}

func TestRunNilSetReturnsError(t *testing.T) {
	p, _ := NewPipeline("p")
	_, err := p.Run(nil)
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestRunAppliesStepsInOrder(t *testing.T) {
	p, _ := NewPipeline("p")
	src := makePipelineSet(t)

	// step 1: uppercase all values
	_ = p.AddStep(func(s *Set) (*Set, error) {
		return Transform(s, func(_, v string) string { return strings.ToUpper(v) })
	})
	// step 2: prefix all keys with APP_
	_ = p.AddStep(func(s *Set) (*Set, error) {
		return CloneWithPrefix(s, s.Name(), "APP_")
	})

	out, err := p.Run(src)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	v, err := out.Get("APP_HOST")
	if err != nil {
		t.Fatalf("Get APP_HOST: %v", err)
	}
	if v != "LOCALHOST" {
		t.Fatalf("expected LOCALHOST, got %q", v)
	}
}

func TestRunDoesNotMutateOriginal(t *testing.T) {
	p, _ := NewPipeline("p")
	src := makePipelineSet(t)
	_ = p.AddStep(func(s *Set) (*Set, error) {
		return Transform(s, func(_, v string) string { return "CHANGED" })
	})
	_, err := p.Run(src)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	v, _ := src.Get("HOST")
	if v != "localhost" {
		t.Fatalf("original set was mutated, got %q", v)
	}
}

func TestRunStepReturnsErrorHaltsExecution(t *testing.T) {
	p, _ := NewPipeline("p")
	src := makePipelineSet(t)
	_ = p.AddStep(func(s *Set) (*Set, error) {
		return nil, fmt.Errorf("step failure")
	})
	_, err := p.Run(src)
	if err == nil {
		t.Fatal("expected error when step fails")
	}
}
