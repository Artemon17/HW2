package persistence

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"pos/internal/application"
	"pos/internal/domain"
	"pos/internal/domain/vo"
)

type studentDTO struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Evaluations []evaluationDTO `json:"evaluations"`
}

type evaluationDTO struct {
	ID          string `json:"id"`
	CriterionID string `json:"criterion_id"`
	Attempt     int    `json:"attempt"`
	Score       int    `json:"score"`
	Comment     string `json:"comment"`
	Date        string `json:"date"`
}

type JSONStudentRepo struct {
	path string
	mu   sync.Mutex
}

func NewJSONStudentRepo(path string) *JSONStudentRepo {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	return &JSONStudentRepo{path: path}
}

func (r *JSONStudentRepo) load() ([]studentDTO, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []studentDTO{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []studentDTO{}, nil
	}
	var dtos []studentDTO
	if err := json.Unmarshal(data, &dtos); err != nil {
		return nil, err
	}
	return dtos, nil
}

func (r *JSONStudentRepo) save(dtos []studentDTO) error {
	data, err := json.MarshalIndent(dtos, "", "  ")
	if err != nil {
		return err
	}
	tmp := r.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, r.path)
}

func (r *JSONStudentRepo) Save(s domain.Student) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	dtos, err := r.load()
	if err != nil {
		return err
	}
	dto := studentToDTO(s)
	for i, d := range dtos {
		if d.ID == dto.ID {
			dtos[i] = dto
			return r.save(dtos)
		}
	}
	return r.save(append(dtos, dto))
}

func (r *JSONStudentRepo) FindByID(id string) (domain.Student, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	dtos, err := r.load()
	if err != nil {
		return domain.Student{}, err
	}
	for _, d := range dtos {
		if d.ID == id {
			return dtoToStudent(d)
		}
	}
	return domain.Student{}, application.ErrStudentNotFound
}

func (r *JSONStudentRepo) FindAll() ([]domain.Student, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	dtos, err := r.load()
	if err != nil {
		return nil, err
	}
	out := make([]domain.Student, 0, len(dtos))
	for _, d := range dtos {
		s, err := dtoToStudent(d)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

func (r *JSONStudentRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	dtos, err := r.load()
	if err != nil {
		return err
	}
	out := dtos[:0]
	found := false
	for _, d := range dtos {
		if d.ID == id {
			found = true
			continue
		}
		out = append(out, d)
	}
	if !found {
		return application.ErrStudentNotFound
	}
	return r.save(out)
}

func studentToDTO(s domain.Student) studentDTO {
	dto := studentDTO{ID: s.ID().String(), Name: s.Name()}
	for _, e := range s.Evaluations() {
		dto.Evaluations = append(dto.Evaluations, evaluationDTO{
			ID:          e.ID().String(),
			CriterionID: e.CriterionID().String(),
			Attempt:     e.Attempt(),
			Score:       e.Score().Value(),
			Comment:     e.Comment(),
			Date:        e.Date().Format(time.RFC3339),
		})
	}
	return dto
}

func dtoToStudent(d studentDTO) (domain.Student, error) {
	sid, err := vo.NewStudentID(d.ID)
	if err != nil {
		return domain.Student{}, err
	}
	evals := make([]domain.Evaluation, 0, len(d.Evaluations))
	for _, ed := range d.Evaluations {
		e, err := dtoToEvaluation(ed)
		if err != nil {
			return domain.Student{}, err
		}
		evals = append(evals, e)
	}
	return domain.ReconstructStudent(sid, d.Name, evals), nil
}

func dtoToEvaluation(ed evaluationDTO) (domain.Evaluation, error) {
	eid, err := vo.NewEvaluationID(ed.ID)
	if err != nil {
		return domain.Evaluation{}, err
	}
	cid, err := vo.NewCriterionID(ed.CriterionID)
	if err != nil {
		return domain.Evaluation{}, err
	}
	score, err := vo.NewScore(ed.Score)
	if err != nil {
		return domain.Evaluation{}, err
	}
	date, err := time.Parse(time.RFC3339, ed.Date)
	if err != nil {
		return domain.Evaluation{}, err
	}
	return domain.ReconstructEvaluation(eid, cid, ed.Attempt, score, ed.Comment, date), nil
}
