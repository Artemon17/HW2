package vo

import "fmt"

type Weight struct{ value float64 }

func NewWeight(v float64) (Weight, error) {
	if v <= 0 {
		return Weight{}, fmt.Errorf("weight must be positive: %v", v)
	}
	return Weight{value: v}, nil
}
func (w Weight) Value() float64   { return w.value }
func (w Weight) Percent() float64 { return w.value * 100 }
