package domain

import (
	"math"
	"pos/internal/domain/vo"
	"testing"
	"time"
)

func buildRubric(t *testing.T, specs ...struct {
	id     string
	max    int
	weight float64
}) Rubric {
	t.Helper()
	r := NewRubric()
	for _, s := range specs {
		if _, err := r.AddCriterion(s.id, mustMax(t, s.max), mustWeight(t, s.weight)); err != nil {
			t.Fatalf("add criterion %q: %v", s.id, err)
		}
	}
	return r
}

func buildStudent(t *testing.T, id string, evals ...struct {
	crit  string
	score int
}) Student {
	t.Helper()
	sid, _ := vo.NewStudentID(id)
	s, _ := NewStudent(sid, "Test")
	for _, e := range evals {
		cid, _ := vo.NewCriterionID(e.crit)
		sc, _ := vo.NewScore(e.score)
		max, _ := vo.NewMaxScore(100) // big enough
		if _, err := s.AddEvaluation(cid, sc, max, "", time.Now()); err != nil {
			t.Fatalf("add eval: %v", err)
		}
	}
	return s
}

func TestCalculator_NoCriteria(t *testing.T) {
	calc := NewFinalGradeCalculator()
	s := buildStudent(t, "s1", struct {
		crit  string
		score int
	}{"hw", 5})
	if _, err := calc.Calculate(s, NewRubric(), LastAttemptPolicy{}); err != ErrNoCriteria {
		t.Errorf("expected ErrNoCriteria, got %v", err)
	}
}

func TestCalculator_NoEvaluations(t *testing.T) {
	calc := NewFinalGradeCalculator()
	r := buildRubric(t, struct {
		id     string
		max    int
		weight float64
	}{"hw", 10, 1.0})
	sid, _ := vo.NewStudentID("s1")
	s, _ := NewStudent(sid, "Empty")
	if _, err := calc.Calculate(s, r, LastAttemptPolicy{}); err != ErrNoEvaluations {
		t.Errorf("expected ErrNoEvaluations, got %v", err)
	}
}

func TestCalculator_LastPolicy(t *testing.T) {
	calc := NewFinalGradeCalculator()
	// one criterion, max=10, weight=1.0
	// two attempts: 5 then 8 -> last = 8 -> 8/10 * 100 = 80
	r := buildRubric(t, struct {
		id     string
		max    int
		weight float64
	}{"hw", 10, 1.0})

	sid, _ := vo.NewStudentID("s1")
	s, _ := NewStudent(sid, "Alice")
	cid, _ := vo.NewCriterionID("hw")
	max, _ := vo.NewMaxScore(10)
	sc1, _ := vo.NewScore(5)
	sc2, _ := vo.NewScore(8)
	s.AddEvaluation(cid, sc1, max, "", time.Now())
	s.AddEvaluation(cid, sc2, max, "", time.Now())

	g, err := calc.Calculate(s, r, LastAttemptPolicy{})
	if err != nil {
		t.Fatalf("calc failed: %v", err)
	}
	if math.Abs(g.Value()-80) > 1e-6 {
		t.Errorf("grade = %v, want 80", g.Value())
	}
}

func TestCalculator_HighestPolicy(t *testing.T) {
	calc := NewFinalGradeCalculator()
	r := buildRubric(t, struct {
		id     string
		max    int
		weight float64
	}{"hw", 10, 1.0})

	sid, _ := vo.NewStudentID("s1")
	s, _ := NewStudent(sid, "Alice")
	cid, _ := vo.NewCriterionID("hw")
	max, _ := vo.NewMaxScore(10)
	sc1, _ := vo.NewScore(9) // best
	sc2, _ := vo.NewScore(4)
	s.AddEvaluation(cid, sc1, max, "", time.Now())
	s.AddEvaluation(cid, sc2, max, "", time.Now())

	g, err := calc.Calculate(s, r, HighestScorePolicy{})
	if err != nil {
		t.Fatalf("calc failed: %v", err)
	}
	if math.Abs(g.Value()-90) > 1e-6 {
		t.Errorf("grade = %v, want 90", g.Value())
	}
}

func TestCalculator_WeightedAverage(t *testing.T) {
	calc := NewFinalGradeCalculator()
	// hw: max=10, w=0.3; exam: max=20, w=0.7
	// student has hw=8 (80%), exam=14 (70%)
	// result = (0.8*0.3 + 0.7*0.7) / (0.3+0.7) * 100
	//        = (0.24 + 0.49) / 1.0 * 100 = 73
	r := buildRubric(t,
		struct {
			id     string
			max    int
			weight float64
		}{"hw", 10, 0.3},
		struct {
			id     string
			max    int
			weight float64
		}{"exam", 20, 0.7},
	)

	sid, _ := vo.NewStudentID("s1")
	s, _ := NewStudent(sid, "Bob")
	maxHW, _ := vo.NewMaxScore(10)
	maxExam, _ := vo.NewMaxScore(20)
	cidHW, _ := vo.NewCriterionID("hw")
	cidExam, _ := vo.NewCriterionID("exam")
	scHW, _ := vo.NewScore(8)
	scExam, _ := vo.NewScore(14)
	s.AddEvaluation(cidHW, scHW, maxHW, "", time.Now())
	s.AddEvaluation(cidExam, scExam, maxExam, "", time.Now())

	g, err := calc.Calculate(s, r, LastAttemptPolicy{})
	if err != nil {
		t.Fatalf("calc failed: %v", err)
	}
	if math.Abs(g.Value()-73) > 1e-6 {
		t.Errorf("grade = %v, want 73", g.Value())
	}
}