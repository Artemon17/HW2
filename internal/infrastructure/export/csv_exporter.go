package export

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"pos/internal/application/dto"
)

type CSVExporter struct{}

func NewCSVExporter() *CSVExporter { return &CSVExporter{} }

func (e *CSVExporter) Format() string { return "csv" }

func (e *CSVExporter) Export(rows []dto.FinalGradeRow) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write([]string{"student_id", "student_name", "final_grade"}); err != nil {
		return nil, err
	}
	for _, r := range rows {
		if err := w.Write([]string{r.StudentID, r.StudentName, fmt.Sprintf("%.2f", r.FinalGrade)}); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return buf.Bytes(), w.Error()
}
