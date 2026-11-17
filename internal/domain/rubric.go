package domain

import (
	"strings"

	"pos/internal/domain/vo"
)

type Rubric struct {
	criteria []Criterion
}

func NewRubric() Rubric { return Rubric{} }

func ReconstructRubric(criteria []Criterion) Rubric {
	return Rubric{criteria: criteria}
}

func (r Rubric) Criteria() []Criterion {
	out := make([]Criterion, len(r.criteria))
	copy(out, r.criteria)
	return out
}

func (r Rubric) FindCriterion(id vo.CriterionID) (Criterion, bool) {
	for _, c := range r.criteria {
		if c.id.Equal(id) {
			return c, true
		}
	}
	return Criterion{}, false
}

func (r *Rubric) AddCriterion(name string, maxScore vo.MaxScore, weight vo.Weight) (vo.CriterionID, error) {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	id, err := vo.NewCriterionID(slug)
	if err != nil {
		return vo.CriterionID{}, ErrEmptyCriterionName
	}
	if _, exists := r.FindCriterion(id); exists {
		return vo.CriterionID{}, ErrDuplicateCriterion
	}
	c, err := NewCriterion(id, name, maxScore, weight)
	if err != nil {
		return vo.CriterionID{}, err
	}
	r.criteria = append(r.criteria, c)
	return id, nil
}

func (r *Rubric) RemoveCriterion(id vo.CriterionID) error {
	for i, c := range r.criteria {
		if c.id.Equal(id) {
			r.criteria = append(r.criteria[:i], r.criteria[i+1:]...)
			return nil
		}
	}
	return ErrCriterionNotFound
}
