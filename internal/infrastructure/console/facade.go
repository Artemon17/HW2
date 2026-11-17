package console

import (
	"fmt"

	"pos/internal/application/usecase"
	"pos/internal/domain"
	"pos/internal/domain/vo"
)

type Deps struct {
	AddStudent          *usecase.AddStudent
	ListStudents        *usecase.ListStudents
	AddCriterion        *usecase.AddCriterion
	ListCriteria        *usecase.ListCriteria
	RemoveCriterion     *usecase.RemoveCriterion
	AddEvaluation       *usecase.AddEvaluation
	ListEvaluations     *usecase.ListEvaluations
	CorrectEvaluation   *usecase.CorrectEvaluation
	RemoveEvaluation    *usecase.RemoveEvaluation
	CalculateFinalGrade *usecase.CalculateFinalGrade
	ExportFinalGrades   *usecase.ExportFinalGrades
}

type Facade struct{ Deps }

func NewFacade(d Deps) *Facade { return &Facade{Deps: d} }

func (f *Facade) printCriteria() {
	criteria, err := f.ListCriteria.Execute()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	if len(criteria) == 0 {
		fmt.Println("Рубрикатор пуст")
		return
	}
	fmt.Println("Критерии:")
	for _, c := range criteria {
		fmt.Printf("  - %s | id=%s | max=%d | weight=%.2f\n",
			c.Name(), c.ID().String(), c.MaxScore().Value(), c.Weight().Value())
	}
}

func (f *Facade) printEvaluations(studentID string) {
	student, err := f.ListEvaluations.Execute(studentID)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	evals := student.Evaluations()
	if len(evals) == 0 {
		fmt.Println("Нет оценок")
		return
	}
	fmt.Printf("Оценки студента %s (%s):\n", student.Name(), student.ID().String())
	for _, e := range evals {
		fmt.Printf("  - eval=%s | crit=%s | attempt=%d | score=%d | %s\n",
			e.ID().String(), e.CriterionID().String(), e.Attempt(), e.Score().Value(), e.Comment())
	}
}

var _ = vo.FinalGrade{}
var _ = domain.LastAttemptPolicy{}
