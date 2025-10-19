package memory

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/dkotsyuruba/go-shortener/internal/model"
)

type MemoryRepository struct {
	filename string
	data     map[string]*model.Link
	mu       sync.RWMutex
}

func NewMemoryRepository(filename string) (*MemoryRepository, error) {
	repo := &MemoryRepository{
		filename: filename,
		data:     make(map[string]*model.Link),
	}

	if err := repo.Init(); err != nil {
		return repo, err
	}

	return repo, nil
}

func (m *MemoryRepository) Init() error {
	info, err := os.Stat(m.filename)
	if os.IsNotExist(err) {
		return err
	}

	if info.Size() == 0 {
		return nil
	}

	content, err := os.ReadFile(m.filename)
	if err != nil {
		return err
	}

	var links map[string]*model.Link
	err = json.Unmarshal(content, &links)
	if err != nil {
		return err
	}

	m.data = links
	return nil
}

func (m *MemoryRepository) Save(link *model.Link) (*model.Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[link.ID]; exists {
		return link, model.ErrDuplicatedURL
	}

	for _, existingLink := range m.data {
		if existingLink.OriginalURL == link.OriginalURL {
			return existingLink, model.ErrDuplicatedURL
		}
	}

	m.data[link.ID] = link
	return link, nil
}

func (m *MemoryRepository) SaveAll(links []*model.Link) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, link := range links {
		if _, exists := m.data[link.ID]; exists {
			return model.ErrDuplicatedURL
		}

		for _, existing := range m.data {
			if existing.OriginalURL == link.OriginalURL {
				return model.ErrDuplicatedURL
			}
		}
	}

	for _, link := range links {
		m.data[link.ID] = link
	}

	return nil
}

func (m *MemoryRepository) FindByID(id string) (*model.Link, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	link, ok := m.data[id]

	return link, ok
}

func (m *MemoryRepository) FindAllByUserID(userID string) ([]*model.Link, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var links []*model.Link

	for _, link := range m.data {
		if link.UUID == userID {
			links = append(links, link)
		}
	}

	if len(links) > 0 {
		return links, nil
	}

	return nil, model.ErrNotFound
}

func (m *MemoryRepository) Close() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	file, err := os.OpenFile(m.filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(m.data)
}

func (m *MemoryRepository) Ping() error {
	_, err := os.Stat(m.filename)
	if os.IsNotExist(err) {
		return err
	}
	return nil
}
