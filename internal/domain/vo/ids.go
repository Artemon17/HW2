package vo

import "fmt"

type StudentID struct{ value string }

func NewStudentID(v string) (StudentID, error) {
	if v == "" {
		return StudentID{}, fmt.Errorf("student id cannot be empty")
	}
	return StudentID{value: v}, nil
}
func (s StudentID) String() string         { return s.value }
func (s StudentID) Equal(o StudentID) bool { return s.value == o.value }

type CriterionID struct{ value string }

func NewCriterionID(v string) (CriterionID, error) {
	if v == "" {
		return CriterionID{}, fmt.Errorf("criterion id cannot be empty")
	}
	return CriterionID{value: v}, nil
}
func (c CriterionID) String() string           { return c.value }
func (c CriterionID) Equal(o CriterionID) bool { return c.value == o.value }

type EvaluationID struct{ value string }

func NewEvaluationID(v string) (EvaluationID, error) {
	if v == "" {
		return EvaluationID{}, fmt.Errorf("evaluation id cannot be empty")
	}
	return EvaluationID{value: v}, nil
}
func (e EvaluationID) String() string            { return e.value }
func (e EvaluationID) Equal(o EvaluationID) bool { return e.value == o.value }
