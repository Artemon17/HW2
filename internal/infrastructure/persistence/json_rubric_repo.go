package persistence

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"pos/internal/domain"
	"pos/internal/domain/vo"
)

type rubricDTO struct {
	Criteria []criterionDTO `json:"criteria"`
}

type criterionDTO struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	MaxScore int     `json:"max_score"`
	Weight   float64 `json:"weight"`
}

type JSONRubricRepo struct {
	path string
	mu   sync.Mutex
}

func NewJSONRubricRepo(path string) *JSONRubricRepo {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	return &JSONRubricRepo{path: path}
}

func (r *JSONRubricRepo) Load() (domain.Rubric, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	data, err := os.ReadFile(r.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain.NewRubric(), nil
		}
		return domain.Rubric{}, err
	}
	if len(data) == 0 {
		return domain.NewRubric(), nil
	}
	var dto rubricDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return domain.Rubric{}, err
	}
	criteria := make([]domain.Criterion, 0, len(dto.Criteria))
	for _, cd := range dto.Criteria {
		c, err := criterionFromDTO(cd)
		if err != nil {
			return domain.Rubric{}, err
		}
		criteria = append(criteria, c)
	}
	return domain.ReconstructRubric(criteria), nil
}

func (r *JSONRubricRepo) Save(rubric domain.Rubric) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	dto := rubricDTO{}
	for _, c := range rubric.Criteria() {
		dto.Criteria = append(dto.Criteria, criterionDTO{
			ID:       c.ID().String(),
			Name:     c.Name(),
			MaxScore: c.MaxScore().Value(),
			Weight:   c.Weight().Value(),
		})
	}
	data, err := json.MarshalIndent(dto, "", "  ")
	if err != nil {
		return err
	}
	tmp := r.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, r.path)
}

func criterionFromDTO(cd criterionDTO) (domain.Criterion, error) {
	id, err := vo.NewCriterionID(cd.ID)
	if err != nil {
		return domain.Criterion{}, err
	}
	maxScore, err := vo.NewMaxScore(cd.MaxScore)
	if err != nil {
		return domain.Criterion{}, err
	}
	weight, err := vo.NewWeight(cd.Weight)
	if err != nil {
		return domain.Criterion{}, err
	}
	return domain.ReconstructCriterion(id, cd.Name, maxScore, weight), nil
}
