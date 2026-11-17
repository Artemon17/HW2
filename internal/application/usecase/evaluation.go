package usecase

import (
	"time"

	"pos/internal/application/port"
	"pos/internal/domain"
	"pos/internal/domain/vo"
)

type AddEvaluation struct {
	studentRepo port.StudentRepository
	rubricRepo  port.RubricRepository
}

func NewAddEvaluation(s port.StudentRepository, r port.RubricRepository) *AddEvaluation {
	return &AddEvaluation{studentRepo: s, rubricRepo: r}
}

func (uc *AddEvaluation) Execute(studentID, criterionID string, score int, comment string) (string, error) {
	sid, err := vo.NewStudentID(studentID)
	if err != nil {
		return "", err
	}
	cid, err := vo.NewCriterionID(criterionID)
	if err != nil {
		return "", err
	}
	sc, err := vo.NewScore(score)
	if err != nil {
		return "", err
	}
	rubric, err := uc.rubricRepo.Load()
	if err != nil {
		return "", err
	}
	criterion, ok := rubric.FindCriterion(cid)
	if !ok {
		return "", domain.ErrCriterionNotFound
	}
	student, err := uc.studentRepo.FindByID(sid.String())
	if err != nil {
		return "", err
	}
	evalID, err := student.AddEvaluation(cid, sc, criterion.MaxScore(), comment, time.Now())
	if err != nil {
		return "", err
	}
	if err := uc.studentRepo.Save(student); err != nil {
		return "", err
	}
	return evalID.String(), nil
}

type ListEvaluations struct{ repo port.StudentRepository }

func NewListEvaluations(repo port.StudentRepository) *ListEvaluations {
	return &ListEvaluations{repo: repo}
}

func (uc *ListEvaluations) Execute(studentID string) (domain.Student, error) {
	sid, err := vo.NewStudentID(studentID)
	if err != nil {
		return domain.Student{}, err
	}
	return uc.repo.FindByID(sid.String())
}

type CorrectEvaluation struct {
	studentRepo port.StudentRepository
	rubricRepo  port.RubricRepository
}

func NewCorrectEvaluation(s port.StudentRepository, r port.RubricRepository) *CorrectEvaluation {
	return &CorrectEvaluation{studentRepo: s, rubricRepo: r}
}

func (uc *CorrectEvaluation) Execute(studentID, evalID string, newScore int) error {
	sid, err := vo.NewStudentID(studentID)
	if err != nil {
		return err
	}
	eid, err := vo.NewEvaluationID(evalID)
	if err != nil {
		return err
	}
	sc, err := vo.NewScore(newScore)
	if err != nil {
		return err
	}
	student, err := uc.studentRepo.FindByID(sid.String())
	if err != nil {
		return err
	}
	var criterionID vo.CriterionID
	found := false
	for _, e := range student.Evaluations() {
		if e.ID().Equal(eid) {
			criterionID = e.CriterionID()
			found = true
			break
		}
	}
	if !found {
		return domain.ErrEvaluationNotFound
	}
	rubric, err := uc.rubricRepo.Load()
	if err != nil {
		return err
	}
	criterion, ok := rubric.FindCriterion(criterionID)
	if !ok {
		return domain.ErrCriterionNotFound
	}
	if err := student.CorrectEvaluation(eid, sc, criterion.MaxScore()); err != nil {
		return err
	}
	return uc.studentRepo.Save(student)
}

type RemoveEvaluation struct{ repo port.StudentRepository }

func NewRemoveEvaluation(repo port.StudentRepository) *RemoveEvaluation {
	return &RemoveEvaluation{repo: repo}
}

func (uc *RemoveEvaluation) Execute(studentID, evalID string) error {
	sid, err := vo.NewStudentID(studentID)
	if err != nil {
		return err
	}
	eid, err := vo.NewEvaluationID(evalID)
	if err != nil {
		return err
	}
	student, err := uc.repo.FindByID(sid.String())
	if err != nil {
		return err
	}
	if err := student.RemoveEvaluation(eid); err != nil {
		return err
	}
	return uc.repo.Save(student)
}
