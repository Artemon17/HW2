package port

import "pos/internal/application/dto"

type GradeExporter interface {
	Format() string
	Export(rows []dto.FinalGradeRow) ([]byte, error)
}
