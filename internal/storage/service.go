package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FaceData struct {
	PersonID  string    `json:"person_id"`
	Embedding []float32 `json:"embedding"`
	CreatedAt int64     `json:"created_at"`
}

type Service interface {
	SaveFace(personID string, embedding []float32) error
	GetAllFaces() ([]FaceData, error)
	GetFaceByID(personID string) (*FaceData, error)
}

type service struct {
	mu    sync.RWMutex
	faces map[string]FaceData
	path  string
}

func NewService() *service {
	s := &service{
		faces: make(map[string]FaceData),
		path:  "data/faces.json",
	}
	s.loadFaces()
	return s
}

func (s *service) SaveFace(personID string, embedding []float32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	face := FaceData{
		PersonID:  personID,
		Embedding: embedding,
		CreatedAt: time.Now().Unix(),
	}

	s.faces[personID] = face
	return s.saveToDisk()
}

func (s *service) GetAllFaces() ([]FaceData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	faces := make([]FaceData, 0, len(s.faces))
	for _, face := range s.faces {
		faces = append(faces, face)
	}
	return faces, nil
}

func (s *service) GetFaceByID(personID string) (*FaceData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if face, exists := s.faces[personID]; exists {
		return &face, nil
	}
	return nil, nil
}

func (s *service) loadFaces() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.faces)
}

func (s *service) saveToDisk() error {
	data, err := json.MarshalIndent(s.faces, "", "  ")
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}
