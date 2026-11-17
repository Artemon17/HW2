package vo

import "testing"

func TestNewScore(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{"zero is valid", 0, false},
		{"positive is valid", 42, false},
		{"negative is invalid", -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewScore(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewScore(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
			if err == nil && s.Value() != tt.value {
				t.Errorf("Value() = %d, want %d", s.Value(), tt.value)
			}
		})
	}
}

func TestScoreIsWithin(t *testing.T) {
	max, _ := NewMaxScore(10)
	cases := []struct {
		score int
		want  bool
	}{
		{0, true},
		{10, true},
		{5, true},
		{11, false},
	}
	for _, c := range cases {
		s, _ := NewScore(c.score)
		if got := s.IsWithin(max); got != c.want {
			t.Errorf("Score(%d).IsWithin(10) = %v, want %v", c.score, got, c.want)
		}
	}
}

func TestNewMaxScore(t *testing.T) {
	if _, err := NewMaxScore(0); err == nil {
		t.Error("MaxScore(0) should be invalid")
	}
	if _, err := NewMaxScore(-5); err == nil {
		t.Error("MaxScore(-5) should be invalid")
	}
	if m, err := NewMaxScore(1); err != nil || m.Value() != 1 {
		t.Errorf("MaxScore(1) unexpected: %v, %v", m, err)
	}
}

func TestNewWeight(t *testing.T) {
	if _, err := NewWeight(0); err == nil {
		t.Error("Weight(0) should be invalid")
	}
	if _, err := NewWeight(-0.1); err == nil {
		t.Error("Weight(-0.1) should be invalid")
	}
	w, err := NewWeight(0.3)
	if err != nil {
		t.Fatalf("NewWeight(0.3) error: %v", err)
	}
	if w.Value() != 0.3 {
		t.Errorf("Value() = %v, want 0.3", w.Value())
	}
	if w.Percent() != 30 {
		t.Errorf("Percent() = %v, want 30", w.Percent())
	}
}

func TestNewFinalGrade(t *testing.T) {
	cases := []struct {
		v       float64
		wantErr bool
	}{
		{0, false},
		{100, false},
		{50.5, false},
		{-0.1, true},
		{100.1, true},
	}
	for _, c := range cases {
		_, err := NewFinalGrade(c.v)
		if (err != nil) != c.wantErr {
			t.Errorf("NewFinalGrade(%v) error = %v, wantErr %v", c.v, err, c.wantErr)
		}
	}
}

func TestIDsCannotBeEmpty(t *testing.T) {
	if _, err := NewStudentID(""); err == nil {
		t.Error("empty StudentID should fail")
	}
	if _, err := NewCriterionID(""); err == nil {
		t.Error("empty CriterionID should fail")
	}
	if _, err := NewEvaluationID(""); err == nil {
		t.Error("empty EvaluationID should fail")
	}
}

func TestIDEquality(t *testing.T) {
	a, _ := NewStudentID("s1")
	b, _ := NewStudentID("s1")
	c, _ := NewStudentID("s2")
	if !a.Equal(b) {
		t.Error("equal IDs should be equal")
	}
	if a.Equal(c) {
		t.Error("different IDs should not be equal")
	}
}
