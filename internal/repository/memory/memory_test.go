package memory

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dkotsyuruba/go-shortener/internal/model"
)

func TestSave(t *testing.T) {
	repo := NewMemoryRepository("test.json")

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
	repo := NewMemoryRepository("test.json")

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
	repo := NewMemoryRepository("test.json")

	link, ok := repo.FindByID("nonexistent-id")
	assert.Nil(t, link)
	assert.False(t, ok)
}

func TestFindByIDSuccess(t *testing.T) {
	repo := NewMemoryRepository("test.json")
	defer os.Remove(repo.filename)

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

func TestPersistSuccess(t *testing.T) {
	repo := NewMemoryRepository("test.json")
	defer os.Remove(repo.filename)

	link := &model.Link{
		ID:          "ghi789",
		OriginalURL: "https://yet-another-site.net",
	}

	err := repo.Save(link)
	require.NoError(t, err)

	err = repo.Persist(repo.filename)
	require.NoError(t, err)

	_, err = os.Stat(repo.filename)
	require.NoError(t, err)
}

func TestLoadFromFile(t *testing.T) {
	t.Run("load from existing file", func(t *testing.T) {
		repo := NewMemoryRepository("test.json")

		tmpFile, _ := os.CreateTemp("", "test.json")
		defer os.Remove(tmpFile.Name())

		jsonData := []byte(`{"jkl012":{"ID":"jkl012","OriginalURL":"https://some-url.com"}}`)
		os.WriteFile(tmpFile.Name(), jsonData, 0644)

		loadedLinks, err := repo.LoadFromFile(tmpFile.Name())
		require.NoError(t, err)
		require.NotNil(t, loadedLinks)
		require.Len(t, loadedLinks, 1)
	})

	t.Run("file not exist", func(t *testing.T) {
		repo := NewMemoryRepository("nonexistent-file.json")
		defer os.Remove(repo.filename)

		_, err := repo.LoadFromFile(repo.filename)
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("empty file", func(t *testing.T) {
		repo := NewMemoryRepository("empty_file.json")
		defer os.Remove(repo.filename)

		f, _ := os.Create(repo.filename)
		f.Close()

		loadedLinks, err := repo.LoadFromFile(repo.filename)
		require.NoError(t, err)
		require.Nil(t, loadedLinks)
	})
}
