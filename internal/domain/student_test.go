package domain

import (
	"pos/internal/domain/vo"
	"testing"
	"time"
)

func mustID(t *testing.T, s string) vo.StudentID {
	t.Helper()
	id, err := vo.NewStudentID(s)
	if err != nil {
		t.Fatalf("bad student id %q: %v", s, err)
	}
	return id
}

func mustCriterionID(t *testing.T, s string) vo.CriterionID {
	t.Helper()
	id, err := vo.NewCriterionID(s)
	if err != nil {
		t.Fatalf("bad criterion id %q: %v", s, err)
	}
	return id
}

func TestNewStudent(t *testing.T) {
	if _, err := NewStudent(mustID(t, "s1"), ""); err == nil {
		t.Error("empty name should fail")
	}
	if _, err := NewStudent(mustID(t, "s1"), "   "); err == nil {
		t.Error("whitespace-only name should fail")
	}
	s, err := NewStudent(mustID(t, "s1"), "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name() != "Alice" {
		t.Errorf("Name() = %q, want %q", s.Name(), "Alice")
	}
}

func TestStudentAddEvaluation(t *testing.T) {
	s, _ := NewStudent(mustID(t, "s1"), "Alice")
	cid := mustCriterionID(t, "hw")
	max, _ := vo.NewMaxScore(10)
	score, _ := vo.NewScore(8)

	id, err := s.AddEvaluation(cid, score, max, "ok", time.Now())
	if err != nil {
		t.Fatalf("AddEvaluation failed: %v", err)
	}
	if id.String() == "" {
		t.Error("evaluation id should not be empty")
	}
	if len(s.Evaluations()) != 1 {
		t.Errorf("expected 1 evaluation, got %d", len(s.Evaluations()))
	}
}

func TestStudentAddEvaluation_ScoreExceedsMax(t *testing.T) {
	s, _ := NewStudent(mustID(t, "s1"), "Alice")
	cid := mustCriterionID(t, "hw")
	max, _ := vo.NewMaxScore(10)
	score, _ := vo.NewScore(11) // > max

	if _, err := s.AddEvaluation(cid, score, max, "", time.Now()); err != ErrScoreExceedsMax {
		t.Errorf("expected ErrScoreExceedsMax, got %v", err)
	}
}

func TestStudentAttemptAutoIncrement(t *testing.T) {
	s, _ := NewStudent(mustID(t, "s1"), "Alice")
	cid := mustCriterionID(t, "hw")
	max, _ := vo.NewMaxScore(10)
	sc1, _ := vo.NewScore(5)
	sc2, _ := vo.NewScore(7)

	s.AddEvaluation(cid, sc1, max, "first", time.Now())
	s.AddEvaluation(cid, sc2, max, "second", time.Now())

	evals := s.Evaluations()
	if len(evals) != 2 {
		t.Fatalf("expected 2 evals, got %d", len(evals))
	}
	if evals[0].Attempt() != 1 || evals[1].Attempt() != 2 {
		t.Errorf("attempts = %d, %d; want 1, 2", evals[0].Attempt(), evals[1].Attempt())
	}
}

func TestStudentCorrectEvaluation(t *testing.T) {
	s, _ := NewStudent(mustID(t, "s1"), "Alice")
	cid := mustCriterionID(t, "hw")
	max, _ := vo.NewMaxScore(10)
	sc, _ := vo.NewScore(5)
	eid, _ := s.AddEvaluation(cid, sc, max, "", time.Now())

	newScore, _ := vo.NewScore(9)
	if err := s.CorrectEvaluation(eid, newScore, max); err != nil {
		t.Fatalf("CorrectEvaluation failed: %v", err)
	}
	if s.Evaluations()[0].Score().Value() != 9 {
		t.Error("score was not updated")
	}
}

func TestStudentCorrectEvaluation_ExceedsMax(t *testing.T) {
	s, _ := NewStudent(mustID(t, "s1"), "Alice")
	cid := mustCriterionID(t, "hw")
	max, _ := vo.NewMaxScore(10)
	sc, _ := vo.NewScore(5)
	eid, _ := s.AddEvaluation(cid, sc, max, "", time.Now())

	badScore, _ := vo.NewScore(15)
	if err := s.CorrectEvaluation(eid, badScore, max); err != ErrScoreExceedsMax {
		t.Errorf("expected ErrScoreExceedsMax, got %v", err)
	}
}

func TestStudentRemoveEvaluation(t *testing.T) {
	s, _ := NewStudent(mustID(t, "s1"), "Alice")
	cid := mustCriterionID(t, "hw")
	max, _ := vo.NewMaxScore(10)
	sc, _ := vo.NewScore(5)
	eid, _ := s.AddEvaluation(cid, sc, max, "", time.Now())

	if err := s.RemoveEvaluation(eid); err != nil {
		t.Fatalf("RemoveEvaluation failed: %v", err)
	}
	if len(s.Evaluations()) != 0 {
		t.Error("evaluation should be removed")
	}
}

func TestStudentRemoveEvaluation_NotFound(t *testing.T) {
	s, _ := NewStudent(mustID(t, "s1"), "Alice")
	badID, _ := vo.NewEvaluationID("nope")
	if err := s.RemoveEvaluation(badID); err != ErrEvaluationNotFound {
		t.Errorf("expected ErrEvaluationNotFound, got %v", err)
	}
}
