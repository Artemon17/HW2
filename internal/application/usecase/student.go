package usecase

import (
	"pos/internal/application"
	"pos/internal/application/port"
	"pos/internal/domain"
	"pos/internal/domain/vo"
)

type AddStudent struct{ repo port.StudentRepository }

func NewAddStudent(repo port.StudentRepository) *AddStudent { return &AddStudent{repo: repo} }

func (uc *AddStudent) Execute(id, name string) error {
	sid, err := vo.NewStudentID(id)
	if err != nil {
		return err
	}
	if _, err := uc.repo.FindByID(sid.String()); err == nil {
		return application.ErrDuplicateStudent
	}
	s, err := domain.NewStudent(sid, name)
	if err != nil {
		return err
	}
	return uc.repo.Save(s)
}

type ListStudents struct{ repo port.StudentRepository }

func NewListStudents(repo port.StudentRepository) *ListStudents { return &ListStudents{repo: repo} }

func (uc *ListStudents) Execute() ([]domain.Student, error) {
	return uc.repo.FindAll()
}
