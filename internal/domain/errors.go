package domain

import "errors"

var (
	ErrEmptyStudentName     = errors.New("student name cannot be empty")
	ErrEmptyCriterionName   = errors.New("criterion name cannot be empty")
	ErrScoreExceedsMax      = errors.New("score exceeds criterion max score")
	ErrEvaluationNotFound   = errors.New("evaluation not found")
	ErrDuplicateCriterion   = errors.New("criterion with this id already exists")
	ErrCriterionNotFound    = errors.New("criterion not found")
	ErrNoCriteria           = errors.New("rubricator has no criteria")
	ErrNoEvaluations        = errors.New("student has no evaluations")
	ErrInconsistentMaxScore = errors.New("criterion max score is below an existing evaluation score")
)
