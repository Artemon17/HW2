package vo

import "fmt"

type Score struct{ value int }

func NewScore(v int) (Score, error) {
	if v < 0 {
		return Score{}, fmt.Errorf("score cannot be negative: %d", v)
	}
	return Score{value: v}, nil
}
func (s Score) Value() int                 { return s.value }
func (s Score) IsWithin(max MaxScore) bool { return s.value <= max.Value() }
