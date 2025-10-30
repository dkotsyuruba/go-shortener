package memory

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/dkotsyuruba/go-shortener/internal/model"
)

type MemoryRepository struct {
	filename string
	data     map[string]*model.Link
	mu       sync.RWMutex
}

func NewMemoryRepository(filename string) *MemoryRepository {
	return &MemoryRepository{
		filename: filename,
		data:     make(map[string]*model.Link),
	}
}

func (m *MemoryRepository) Save(link *model.Link) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[link.ID]; exists {
		return errors.New("duplicate key")
	}

	m.data[link.ID] = link

	return nil
}

func (m *MemoryRepository) FindByID(id string) (*model.Link, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	link, ok := m.data[id]

	return link, ok
}

func (m *MemoryRepository) Persist(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	urls := make([]*model.Link, 0, len(m.data))
	for _, v := range m.data {
		urls = append(urls, v)
	}

	data, err := json.Marshal(urls)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (m *MemoryRepository) LoadFromFile(filename string) (map[string]*model.Link, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil, err
	}

	if info.Size() == 0 {
		return nil, nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var links []*model.Link
	if err := json.Unmarshal(content, &links); err != nil {
		return nil, err
	}

	for _, link := range links {
		m.data[link.ID] = link
	}

	return m.data, nil
}
