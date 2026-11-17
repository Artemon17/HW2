package vo

import "fmt"

type MaxScore struct{ value int }

func NewMaxScore(v int) (MaxScore, error) {
	if v <= 0 {
		return MaxScore{}, fmt.Errorf("max score must be positive: %d", v)
	}
	return MaxScore{value: v}, nil
}
func (m MaxScore) Value() int { return m.value }
