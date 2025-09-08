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
		Id:          "abc123",
		OriginalUrl: "https://example.com",
	}

	err := repo.Save(link)
	require.NoError(t, err)

	foundLink, ok := repo.FindById("abc123")
	require.True(t, ok)
	assert.Equal(t, link, foundLink)
}

func TestSaveDuplicateKey(t *testing.T) {
	repo := NewMemoryRepository()

	link := &model.Link{
		Id:          "abc123",
		OriginalUrl: "https://example.com",
	}

	err := repo.Save(link)
	require.NoError(t, err)

	err = repo.Save(link)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
}

func TestFindByIdNotFound(t *testing.T) {
	repo := NewMemoryRepository()

	link, ok := repo.FindById("nonexistent-id")
	assert.Nil(t, link)
	assert.False(t, ok)
}

func TestFindByIdSuccess(t *testing.T) {
	repo := NewMemoryRepository()

	link := &model.Link{
		Id:          "def456",
		OriginalUrl: "https://another-site.org",
	}

	err := repo.Save(link)
	require.NoError(t, err)

	foundLink, ok := repo.FindById("def456")
	require.True(t, ok)
	assert.Equal(t, link, foundLink)
}
