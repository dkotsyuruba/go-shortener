package memory

import (
	"errors"
	"sync"

	"github.com/dkotsyuruba/go-shortener/internal/model"
)

type MemoryRepository struct {
	data map[string]*model.Link
	mu   sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		data: make(map[string]*model.Link),
	}
}

func (m *MemoryRepository) Save(link *model.Link) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[link.Id]; exists {
		return errors.New("duplicate key")
	}

	m.data[link.Id] = link

	return nil
}

func (m *MemoryRepository) FindById(id string) (*model.Link, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	link, ok := m.data[id]

	return link, ok
}
