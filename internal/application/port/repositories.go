package port

import "pos/internal/domain"

type StudentRepository interface {
	Save(s domain.Student) error
	FindByID(id string) (domain.Student, error)
	FindAll() ([]domain.Student, error)
	Delete(id string) error
}

type RubricRepository interface {
	Load() (domain.Rubric, error)
	Save(r domain.Rubric) error
}
