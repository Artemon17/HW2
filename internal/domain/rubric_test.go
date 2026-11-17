package domain

import (
	"pos/internal/domain/vo"
	"testing"
)

func mustMax(t *testing.T, v int) vo.MaxScore {
	t.Helper()
	m, err := vo.NewMaxScore(v)
	if err != nil {
		t.Fatalf("bad max score: %v", err)
	}
	return m
}

func mustWeight(t *testing.T, v float64) vo.Weight {
	t.Helper()
	w, err := vo.NewWeight(v)
	if err != nil {
		t.Fatalf("bad weight: %v", err)
	}
	return w
}

func TestRubricAddCriterion(t *testing.T) {
	r := NewRubric()
	id, err := r.AddCriterion("Homework", mustMax(t, 10), mustWeight(t, 0.3))
	if err != nil {
		t.Fatalf("AddCriterion failed: %v", err)
	}
	if id.String() != "homework" {
		t.Errorf("slug = %q, want %q", id.String(), "homework")
	}
	if len(r.Criteria()) != 1 {
		t.Error("expected 1 criterion")
	}
}

func TestRubricAddCriterion_EmptyName(t *testing.T) {
	r := NewRubric()
	if _, err := r.AddCriterion("", mustMax(t, 10), mustWeight(t, 0.3)); err != ErrEmptyCriterionName {
		t.Errorf("expected ErrEmptyCriterionName, got %v", err)
	}
}

func TestRubricAddCriterion_Duplicate(t *testing.T) {
	r := NewRubric()
	r.AddCriterion("Homework", mustMax(t, 10), mustWeight(t, 0.3))
	if _, err := r.AddCriterion("Homework", mustMax(t, 5), mustWeight(t, 0.2)); err != ErrDuplicateCriterion {
		t.Errorf("expected ErrDuplicateCriterion, got %v", err)
	}
}

func TestRubricFindCriterion(t *testing.T) {
	r := NewRubric()
	id, _ := r.AddCriterion("Exam", mustMax(t, 20), mustWeight(t, 0.5))
	c, ok := r.FindCriterion(id)
	if !ok {
		t.Fatal("criterion not found")
	}
	if c.Name() != "Exam" {
		t.Errorf("Name() = %q, want %q", c.Name(), "Exam")
	}
	missing, _ := vo.NewCriterionID("missing")
	if _, ok := r.FindCriterion(missing); ok {
		t.Error("should not find missing criterion")
	}
}

func TestRubricRemoveCriterion(t *testing.T) {
	r := NewRubric()
	id, _ := r.AddCriterion("X", mustMax(t, 10), mustWeight(t, 0.5))
	if err := r.RemoveCriterion(id); err != nil {
		t.Fatalf("RemoveCriterion failed: %v", err)
	}
	if len(r.Criteria()) != 0 {
		t.Error("criterion should be removed")
	}
}

func TestRubricRemoveCriterion_NotFound(t *testing.T) {
	r := NewRubric()
	bad, _ := vo.NewCriterionID("nope")
	if err := r.RemoveCriterion(bad); err != ErrCriterionNotFound {
		t.Errorf("expected ErrCriterionNotFound, got %v", err)
	}
}