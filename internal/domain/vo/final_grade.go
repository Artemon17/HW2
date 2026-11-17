package vo

import "fmt"

type FinalGrade struct{ value float64 }

func NewFinalGrade(v float64) (FinalGrade, error) {
	if v < 0 || v > 100 {
		return FinalGrade{}, fmt.Errorf("final grade must be in [0, 100]: %v", v)
	}
	return FinalGrade{value: v}, nil
}
func (f FinalGrade) Value() float64 { return f.value }
func (f FinalGrade) String() string { return fmt.Sprintf("%.2f", f.value) }
