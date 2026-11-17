package usecase

import (
	"fmt"

	"pos/internal/application/dto"
	"pos/internal/application/port"
	"pos/internal/domain"
	"pos/internal/domain/vo"
)

type CalculateFinalGrade struct {
	studentRepo port.StudentRepository
	rubricRepo  port.RubricRepository
	calc        domain.FinalGradeCalculator
}

func NewCalculateFinalGrade(s port.StudentRepository, r port.RubricRepository, c domain.FinalGradeCalculator) *CalculateFinalGrade {
	return &CalculateFinalGrade{studentRepo: s, rubricRepo: r, calc: c}
}

func (uc *CalculateFinalGrade) Execute(studentID string, policy domain.AttemptPolicy) (vo.FinalGrade, error) {
	sid, err := vo.NewStudentID(studentID)
	if err != nil {
		return vo.FinalGrade{}, err
	}
	student, err := uc.studentRepo.FindByID(sid.String())
	if err != nil {
		return vo.FinalGrade{}, err
	}
	rubric, err := uc.rubricRepo.Load()
	if err != nil {
		return vo.FinalGrade{}, err
	}
	return uc.calc.Calculate(student, rubric, policy)
}

type ExportFinalGrades struct {
	studentRepo port.StudentRepository
	rubricRepo  port.RubricRepository
	calc        domain.FinalGradeCalculator
	exporter    port.GradeExporter
}

func NewExportFinalGrades(s port.StudentRepository, r port.RubricRepository, c domain.FinalGradeCalculator, e port.GradeExporter) *ExportFinalGrades {
	return &ExportFinalGrades{studentRepo: s, rubricRepo: r, calc: c, exporter: e}
}

func (uc *ExportFinalGrades) Execute(policy domain.AttemptPolicy) ([]byte, error) {
	students, err := uc.studentRepo.FindAll()
	if err != nil {
		return nil, err
	}
	rubric, err := uc.rubricRepo.Load()
	if err != nil {
		return nil, err
	}
	var rows []dto.FinalGradeRow
	for _, s := range students {
		grade, err := uc.calc.Calculate(s, rubric, policy)
		if err != nil {
			return nil, fmt.Errorf("student %s: %w", s.ID().String(), err)
		}
		rows = append(rows, dto.FinalGradeRow{
			StudentID:   s.ID().String(),
			StudentName: s.Name(),
			FinalGrade:  grade.Value(),
		})
	}
	return uc.exporter.Export(rows)
}
