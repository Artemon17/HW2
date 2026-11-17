package domain

import "pos/internal/domain/vo"

type FinalGradeCalculator struct{}

func NewFinalGradeCalculator() FinalGradeCalculator { return FinalGradeCalculator{} }

func (FinalGradeCalculator) Calculate(student Student, rubric Rubric, policy AttemptPolicy) (vo.FinalGrade, error) {
	criteria := rubric.Criteria()
	if len(criteria) == 0 {
		return vo.FinalGrade{}, ErrNoCriteria
	}
	if len(student.Evaluations()) == 0 {
		return vo.FinalGrade{}, ErrNoEvaluations
	}

	byCriterion := make(map[string][]Evaluation)
	for _, e := range student.Evaluations() {
		key := e.CriterionID().String()
		byCriterion[key] = append(byCriterion[key], e)
	}

	var numerator, totalWeight float64
	for _, c := range criteria {
		var score float64
		if e, ok := policy.Select(byCriterion[c.ID().String()]); ok {
			score = float64(e.Score().Value())
		}
		maxScore := float64(c.MaxScore().Value())
		if score > maxScore {
			return vo.FinalGrade{}, ErrInconsistentMaxScore
		}
		w := c.Weight().Value()
		numerator += (score / maxScore) * w
		totalWeight += w
	}

	if totalWeight == 0 {
		return vo.FinalGrade{}, ErrNoCriteria
	}
	return vo.NewFinalGrade((numerator / totalWeight) * 100)
}
