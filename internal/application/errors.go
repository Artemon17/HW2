package application

import "errors"

var (
	ErrStudentNotFound         = errors.New("student not found")
	ErrDuplicateStudent        = errors.New("student with this id already exists")
	ErrCriterionHasEvaluations = errors.New("cannot remove criterion: it has evaluations")
)
