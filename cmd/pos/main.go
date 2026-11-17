package main

import (
	"fmt"
	"os"
	"path/filepath"

	"pos/internal/application/usecase"
	"pos/internal/domain"
	"pos/internal/infrastructure/console"
	"pos/internal/infrastructure/export"
	"pos/internal/infrastructure/persistence"
)

func main() {
	dataDir := "./data"
	if len(os.Args) > 1 {
		dataDir = os.Args[1]
	}

	studentRepo := persistence.NewJSONStudentRepo(filepath.Join(dataDir, "students.json"))
	rubricRepo := persistence.NewJSONRubricRepo(filepath.Join(dataDir, "rubric.json"))
	exporter := export.NewCSVExporter()
	calc := domain.NewFinalGradeCalculator()

	deps := console.Deps{
		AddStudent:          usecase.NewAddStudent(studentRepo),
		ListStudents:        usecase.NewListStudents(studentRepo),
		AddCriterion:        usecase.NewAddCriterion(rubricRepo),
		ListCriteria:        usecase.NewListCriteria(rubricRepo),
		RemoveCriterion:     usecase.NewRemoveCriterion(rubricRepo, studentRepo),
		AddEvaluation:       usecase.NewAddEvaluation(studentRepo, rubricRepo),
		ListEvaluations:     usecase.NewListEvaluations(studentRepo),
		CorrectEvaluation:   usecase.NewCorrectEvaluation(studentRepo, rubricRepo),
		RemoveEvaluation:    usecase.NewRemoveEvaluation(studentRepo),
		CalculateFinalGrade: usecase.NewCalculateFinalGrade(studentRepo, rubricRepo, calc),
		ExportFinalGrades:   usecase.NewExportFinalGrades(studentRepo, rubricRepo, calc, exporter),
	}

	fmt.Println("Данные:", dataDir)
	console.NewFacade(deps).Run()
}
