package domain

import (
	"pos/internal/domain/vo"
	"time"
)

type Evaluation struct {
	id          vo.EvaluationID
	criterionID vo.CriterionID
	attempt     int
	score       vo.Score
	comment     string
	date        time.Time
}

func newEvaluation(
	id vo.EvaluationID,
	criterionID vo.CriterionID,
	attempt int,
	score vo.Score,
	maxScore vo.MaxScore,
	comment string,
	date time.Time,
) (Evaluation, error) {
	if !score.IsWithin(maxScore) {
		return Evaluation{}, ErrScoreExceedsMax
	}
	return Evaluation{
		id:          id,
		criterionID: criterionID,
		attempt:     attempt,
		score:       score,
		comment:     comment,
		date:        date,
	}, nil
}

func ReconstructEvaluation(
	id vo.EvaluationID,
	criterionID vo.CriterionID,
	attempt int,
	score vo.Score,
	comment string,
	date time.Time,
) Evaluation {
	return Evaluation{
		id:          id,
		criterionID: criterionID,
		attempt:     attempt,
		score:       score,
		comment:     comment,
		date:        date,
	}
}

func (e Evaluation) ID() vo.EvaluationID         { return e.id }
func (e Evaluation) CriterionID() vo.CriterionID { return e.criterionID }
func (e Evaluation) Attempt() int                { return e.attempt }
func (e Evaluation) Score() vo.Score             { return e.score }
func (e Evaluation) Comment() string             { return e.comment }
func (e Evaluation) Date() time.Time             { return e.date }
