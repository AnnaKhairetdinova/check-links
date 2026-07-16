package storage

import (
	"check-links/models"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type JSONStorage struct {
	mu   sync.RWMutex
	path string
	data []models.LinkList
}

func NewJSONStorage(path string) (*JSONStorage, error) {
	s := &JSONStorage{path: path}

	if err := s.load(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *JSONStorage) load() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	file, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&s.data)
}

func (s *JSONStorage) save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(s.data)
}

func (s *JSONStorage) Create(l models.LinkList) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	l.Num = len(s.data) + 1
	l.Status = models.LinkListStatusInProgress

	s.data = append(s.data, l)
	return s.save()
}

func (s *JSONStorage) List() []models.LinkList {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cp := make([]models.LinkList, len(s.data))
	copy(cp, s.data)
	return cp
}

func (s *JSONStorage) UpdateStatus(num int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.data {
		if s.data[i].Num == num {
			if s.data[i].Status == models.LinkListStatusInProgress {
				s.data[i].Status = models.LinkListStatusDone
			}
			return s.save()
		}
	}

	return fmt.Errorf("could not find link %d", num)
}
