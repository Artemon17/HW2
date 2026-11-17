package domain

type AttemptPolicy interface {
	Name() string
	Select(evals []Evaluation) (Evaluation, bool)
}

type LastAttemptPolicy struct{}

func (LastAttemptPolicy) Name() string { return "last" }

func (LastAttemptPolicy) Select(evals []Evaluation) (Evaluation, bool) {
	if len(evals) == 0 {
		return Evaluation{}, false
	}
	best := evals[0]
	for _, e := range evals[1:] {
		if e.Attempt() > best.Attempt() {
			best = e
		}
	}
	return best, true
}

type HighestScorePolicy struct{}

func (HighestScorePolicy) Name() string { return "highest" }

func (HighestScorePolicy) Select(evals []Evaluation) (Evaluation, bool) {
	if len(evals) == 0 {
		return Evaluation{}, false
	}
	best := evals[0]
	for _, e := range evals[1:] {
		if e.Score().Value() > best.Score().Value() {
			best = e
		}
	}
	return best, true
}
