package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndValidateToken(t *testing.T) {
	manager := NewJWTManager("mysecretkey1234567890")

	userID := "aabbccddeeff00112233445566778899"
	token, err := manager.GenerateToken(userID)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	extractedUserID, err := manager.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)
}

func TestGenerateTokenEmptyUserID(t *testing.T) {
	manager := NewJWTManager("mysecretkey1234567890")

	token, err := manager.GenerateToken("")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	_, err = manager.ValidateToken(token)
	require.NoError(t, err)
}

func TestValidateTokenInvalidToken(t *testing.T) {
	manager := NewJWTManager("mysecretkey1234567890")

	_, err := manager.ValidateToken("")
	assert.Error(t, err)

	_, err = manager.ValidateToken("1234")
	assert.Error(t, err)

	badToken := "aabbccddeeff00112233445566778899" + "ffffffffffffffffffffffffffffffffffffffffffffffff"
	_, err = manager.ValidateToken(badToken)
	assert.Error(t, err)
}
