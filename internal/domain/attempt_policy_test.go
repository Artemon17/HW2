package domain

import (
	"pos/internal/domain/vo"
	"testing"
	"time"
)

func mkEval(t *testing.T, attempt, score int) Evaluation {
	t.Helper()
	eid, _ := vo.NewEvaluationID("e")
	cid, _ := vo.NewCriterionID("c")
	sc, _ := vo.NewScore(score)
	max, _ := vo.NewMaxScore(100)
	e, err := newEvaluation(eid, cid, attempt, sc, max, "", time.Now())
	if err != nil {
		t.Fatalf("mkEval: %v", err)
	}
	return e
}

func TestLastAttemptPolicy(t *testing.T) {
	evals := []Evaluation{
		mkEval(t, 1, 5),
		mkEval(t, 3, 2), // last by attempt
		mkEval(t, 2, 9),
	}
	p := LastAttemptPolicy{}
	got, ok := p.Select(evals)
	if !ok {
		t.Fatal("expected selection")
	}
	if got.Attempt() != 3 {
		t.Errorf("attempt = %d, want 3", got.Attempt())
	}
}

func TestHighestScorePolicy(t *testing.T) {
	evals := []Evaluation{
		mkEval(t, 1, 5),
		mkEval(t, 2, 9), // highest
		mkEval(t, 3, 7),
	}
	p := HighestScorePolicy{}
	got, ok := p.Select(evals)
	if !ok {
		t.Fatal("expected selection")
	}
	if got.Score().Value() != 9 {
		t.Errorf("score = %d, want 9", got.Score().Value())
	}
}

func TestPolicyEmpty(t *testing.T) {
	if _, ok := (LastAttemptPolicy{}).Select(nil); ok {
		t.Error("last policy should return false on empty")
	}
	if _, ok := (HighestScorePolicy{}).Select(nil); ok {
		t.Error("highest policy should return false on empty")
	}
}