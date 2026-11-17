package domain

import (
	"strings"

	"pos/internal/domain/vo"
)

type Criterion struct {
	id       vo.CriterionID
	name     string
	maxScore vo.MaxScore
	weight   vo.Weight
}

func NewCriterion(id vo.CriterionID, name string, maxScore vo.MaxScore, weight vo.Weight) (Criterion, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Criterion{}, ErrEmptyCriterionName
	}
	return Criterion{id: id, name: name, maxScore: maxScore, weight: weight}, nil
}

func ReconstructCriterion(id vo.CriterionID, name string, maxScore vo.MaxScore, weight vo.Weight) Criterion {
	return Criterion{id: id, name: name, maxScore: maxScore, weight: weight}
}

func (c Criterion) ID() vo.CriterionID    { return c.id }
func (c Criterion) Name() string          { return c.name }
func (c Criterion) MaxScore() vo.MaxScore { return c.maxScore }
func (c Criterion) Weight() vo.Weight     { return c.weight }
