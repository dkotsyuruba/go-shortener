package memory

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dkotsyuruba/go-shortener/internal/model"
)

var testFilename string

func init() {
	testFilename = createTemporaryJSONFile()
}

func createTemporaryJSONFile() string {
	tmpFile, err := os.CreateTemp("", "test.json")
	if err != nil {
		panic(err)
	}
	defer tmpFile.Close()

	return tmpFile.Name()
}

func TestSaveSuccess(t *testing.T) {
	repo, err := NewMemoryRepository(testFilename)
	require.NoError(t, err)

	link := &model.Link{
		ID:          "abc123",
		OriginalURL: "https://example.com",
	}

	result, err := repo.Save(link)
	require.NoError(t, err)

	assert.Equal(t, link, result)
}

func TestSaveDuplicateLink(t *testing.T) {
	repo, err := NewMemoryRepository(testFilename)
	require.NoError(t, err)

	link := &model.Link{
		ID:          "abc123",
		OriginalURL: "https://example.com",
	}

	_, err = repo.Save(link)
	require.NoError(t, err)

	_, err = repo.Save(link)
	require.Error(t, err)
	assert.EqualError(t, err, model.ErrDuplicatedURL.Error())
}

func TestSaveAllSuccess(t *testing.T) {
	repo, err := NewMemoryRepository(testFilename)
	require.NoError(t, err)

	links := []*model.Link{
		{ID: "ghi789", OriginalURL: "https://site-a.com"},
		{ID: "jkl012", OriginalURL: "https://site-b.com"},
	}

	err = repo.SaveAll(links)
	require.NoError(t, err)

	foundLinkA, okA := repo.FindByID("ghi789")
	require.True(t, okA)
	assert.Equal(t, links[0], foundLinkA)

	foundLinkB, okB := repo.FindByID("jkl012")
	require.True(t, okB)
	assert.Equal(t, links[1], foundLinkB)
}

func TestCloseSuccess(t *testing.T) {
	repo, err := NewMemoryRepository(testFilename)
	require.NoError(t, err)

	err = repo.Close()
	require.NoError(t, err)
}

func TestPingSuccess(t *testing.T) {
	repo, err := NewMemoryRepository(testFilename)
	require.NoError(t, err)

	err = repo.Ping()
	require.NoError(t, err)
}

func TestFindByOriginalURLSuccess(t *testing.T) {
	repo, err := NewMemoryRepository(testFilename)
	require.NoError(t, err)

	link := &model.Link{
		ID:          "def456",
		OriginalURL: "https://another-site.org",
	}

	_, err = repo.Save(link)
	require.NoError(t, err)

	foundLink, ok := repo.FindByOriginalURL("https://another-site.org")
	require.True(t, ok)
	assert.Equal(t, link, foundLink)
}

func TestFindByOriginalURLNotFound(t *testing.T) {
	repo, err := NewMemoryRepository(testFilename)
	require.NoError(t, err)

	link, ok := repo.FindByOriginalURL("https://nonexistent-url.com")
	assert.Nil(t, link)
	assert.False(t, ok)
}
