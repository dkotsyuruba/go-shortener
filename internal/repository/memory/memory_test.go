package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dkotsyuruba/go-shortener/internal/model"
)

func TestSave(t *testing.T) {
	repo := NewMemoryRepository()

	link := &model.Link{
		ID:          "abc123",
		OriginalURL: "https://example.com",
	}

	err := repo.Save(link)
	require.NoError(t, err)

	foundLink, ok := repo.FindByID("abc123")
	require.True(t, ok)
	assert.Equal(t, link, foundLink)
}

func TestSaveDuplicateKey(t *testing.T) {
	repo := NewMemoryRepository()

	link := &model.Link{
		ID:          "abc123",
		OriginalURL: "https://example.com",
	}

	err := repo.Save(link)
	require.NoError(t, err)

	err = repo.Save(link)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
}

func TestFindByIDNotFound(t *testing.T) {
	repo := NewMemoryRepository()

	link, ok := repo.FindByID("nonexistent-id")
	assert.Nil(t, link)
	assert.False(t, ok)
}

func TestFindByIDSuccess(t *testing.T) {
	repo := NewMemoryRepository()

	link := &model.Link{
		ID:          "def456",
		OriginalURL: "https://another-site.org",
	}

	err := repo.Save(link)
	require.NoError(t, err)

	foundLink, ok := repo.FindByID("def456")
	require.True(t, ok)
	assert.Equal(t, link, foundLink)
}
