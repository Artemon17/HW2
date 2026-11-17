package usecase

import (
	"pos/internal/application"
	"pos/internal/application/port"
	"pos/internal/domain"
	"pos/internal/domain/vo"
)

type AddCriterion struct{ repo port.RubricRepository }

func NewAddCriterion(repo port.RubricRepository) *AddCriterion { return &AddCriterion{repo: repo} }

func (uc *AddCriterion) Execute(name string, maxScore int, weight float64) (string, error) {
	ms, err := vo.NewMaxScore(maxScore)
	if err != nil {
		return "", err
	}
	w, err := vo.NewWeight(weight)
	if err != nil {
		return "", err
	}
	rubric, err := uc.repo.Load()
	if err != nil {
		return "", err
	}
	id, err := rubric.AddCriterion(name, ms, w)
	if err != nil {
		return "", err
	}
	if err := uc.repo.Save(rubric); err != nil {
		return "", err
	}
	return id.String(), nil
}

type ListCriteria struct{ repo port.RubricRepository }

func NewListCriteria(repo port.RubricRepository) *ListCriteria { return &ListCriteria{repo: repo} }

func (uc *ListCriteria) Execute() ([]domain.Criterion, error) {
	rubric, err := uc.repo.Load()
	if err != nil {
		return nil, err
	}
	return rubric.Criteria(), nil
}

type RemoveCriterion struct {
	rubricRepo  port.RubricRepository
	studentRepo port.StudentRepository
}

func NewRemoveCriterion(r port.RubricRepository, s port.StudentRepository) *RemoveCriterion {
	return &RemoveCriterion{rubricRepo: r, studentRepo: s}
}

func (uc *RemoveCriterion) Execute(criterionID string) error {
	cid, err := vo.NewCriterionID(criterionID)
	if err != nil {
		return err
	}
	// cross-aggregate rule: cannot remove criterion that has evaluations
	students, err := uc.studentRepo.FindAll()
	if err != nil {
		return err
	}
	for _, s := range students {
		for _, e := range s.Evaluations() {
			if e.CriterionID().Equal(cid) {
				return application.ErrCriterionHasEvaluations
			}
		}
	}
	rubric, err := uc.rubricRepo.Load()
	if err != nil {
		return err
	}
	if err := rubric.RemoveCriterion(cid); err != nil {
		return err
	}
	return uc.rubricRepo.Save(rubric)
}
