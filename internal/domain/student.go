package domain

import (
	"fmt"
	"strings"
	"time"

	"pos/internal/domain/vo"
)

type Student struct {
	id          vo.StudentID
	name        string
	evaluations []Evaluation
}

func NewStudent(id vo.StudentID, name string) (Student, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Student{}, ErrEmptyStudentName
	}
	return Student{id: id, name: name}, nil
}

// ReconstructStudent is for infrastructure (persistence) use only.
func ReconstructStudent(id vo.StudentID, name string, evals []Evaluation) Student {
	return Student{id: id, name: name, evaluations: evals}
}

func (s Student) ID() vo.StudentID { return s.id }
func (s Student) Name() string     { return s.name }

func (s Student) Evaluations() []Evaluation {
	out := make([]Evaluation, len(s.evaluations))
	copy(out, s.evaluations)
	return out
}

func (s *Student) AddEvaluation(
	criterionID vo.CriterionID,
	score vo.Score,
	maxScore vo.MaxScore,
	comment string,
	date time.Time,
) (vo.EvaluationID, error) {
	attempt := 1
	for _, e := range s.evaluations {
		if e.criterionID.Equal(criterionID) && e.attempt >= attempt {
			attempt = e.attempt + 1
		}
	}
	id, err := vo.NewEvaluationID(fmt.Sprintf("eval-%s-%s-%d", s.id.String(), criterionID.String(), attempt))
	if err != nil {
		return vo.EvaluationID{}, err
	}
	e, err := newEvaluation(id, criterionID, attempt, score, maxScore, comment, date)
	if err != nil {
		return vo.EvaluationID{}, err
	}
	s.evaluations = append(s.evaluations, e)
	return id, nil
}

func (s *Student) CorrectEvaluation(evalID vo.EvaluationID, newScore vo.Score, maxScore vo.MaxScore) error {
	for i := range s.evaluations {
		if s.evaluations[i].id.Equal(evalID) {
			if !newScore.IsWithin(maxScore) {
				return ErrScoreExceedsMax
			}
			s.evaluations[i].score = newScore
			return nil
		}
	}
	return ErrEvaluationNotFound
}

func (s *Student) RemoveEvaluation(evalID vo.EvaluationID) error {
	for i := range s.evaluations {
		if s.evaluations[i].id.Equal(evalID) {
			s.evaluations = append(s.evaluations[:i], s.evaluations[i+1:]...)
			return nil
		}
	}
	return ErrEvaluationNotFound
}
