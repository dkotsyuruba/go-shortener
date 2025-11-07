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

func (m *MemoryRepository) Migrate() error {
	return nil
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

	var links []*model.Link
	if err := json.Unmarshal(content, &links); err != nil {
		return err
	}

	for _, link := range links {
		m.data[link.ID] = link
	}

	return nil
}

func (m *MemoryRepository) Save(link *model.Link) (*model.Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[link.ID]; exists {
		return link, model.ErrDuplicatedURL
	}

	foundLink, ok := m.FindByOriginalURL(link.OriginalURL)
	if ok {
		return foundLink, model.ErrDuplicatedURL
	}

	m.data[link.ID] = link

	return link, nil
}

func (m *MemoryRepository) SaveAll(links []*model.Link) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, link := range links {
		if _, exists := m.data[link.ID]; exists {
			continue
		}
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

func (m *MemoryRepository) FindByOriginalURL(url string) (*model.Link, bool) {
	for _, link := range m.data {
		if link.OriginalURL == url {
			return link, true
		}
	}
	return nil, false
}

func (m *MemoryRepository) Close() error {
	file, err := os.OpenFile(m.filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
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

func (m *MemoryRepository) Ping() error {
	_, err := os.Stat(m.filename)
	if os.IsNotExist(err) {
		return err
	}
	return nil
}
