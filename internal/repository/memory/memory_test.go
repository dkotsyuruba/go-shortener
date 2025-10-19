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

func TestNewMemoryRepositoryInit(t *testing.T) {
	repo, err := NewMemoryRepository(testFilename)
	require.NoError(t, err)
	assert.NotNil(t, repo)
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

func TestSaveAllDuplicate(t *testing.T) {
	repo, err := NewMemoryRepository(testFilename)
	require.NoError(t, err)

	existingLink := &model.Link{
		ID:          "dup123",
		OriginalURL: "https://duplicate.com",
	}
	_, err = repo.Save(existingLink)
	require.NoError(t, err)

	links := []*model.Link{
		{ID: "dup123", OriginalURL: "https://site-a.com"},    // duplicate ID
		{ID: "new456", OriginalURL: "https://duplicate.com"}, // duplicate OriginalURL
	}

	err = repo.SaveAll(links)
	require.Error(t, err)
	assert.EqualError(t, err, model.ErrDuplicatedURL.Error())
}

func TestFindByID(t *testing.T) {
	repo, err := NewMemoryRepository(testFilename)
	require.NoError(t, err)

	link := &model.Link{
		ID:          "find123",
		OriginalURL: "https://findme.com",
	}
	_, err = repo.Save(link)
	require.NoError(t, err)

	foundLink, ok := repo.FindByID("find123")
	require.True(t, ok)
	assert.Equal(t, link, foundLink)

	_, ok = repo.FindByID("nonexistent")
	assert.False(t, ok)
}

func TestFindAllByUserID(t *testing.T) {
	repo, err := NewMemoryRepository(testFilename)
	require.NoError(t, err)

	userID := "user-1"
	otherUserID := "user-2"

	links := []*model.Link{
		{ID: "id1", OriginalURL: "https://a.com", UUID: userID},
		{ID: "id2", OriginalURL: "https://b.com", UUID: userID},
		{ID: "id3", OriginalURL: "https://c.com", UUID: otherUserID},
	}

	err = repo.SaveAll(links)
	require.NoError(t, err)

	foundLinks, err := repo.FindAllByUserID(userID)
	require.NoError(t, err)
	assert.Len(t, foundLinks, 2)
	for _, link := range foundLinks {
		assert.Equal(t, userID, link.UUID)
	}

	_, err = repo.FindAllByUserID("nonexistent-user")
	assert.Error(t, err)
	assert.Equal(t, model.ErrNotFound, err)
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

func TestMain(m *testing.M) {
	code := m.Run()
	os.Remove(testFilename)
	os.Exit(code)
}
