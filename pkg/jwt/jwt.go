package jwt

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrMissingUserID = errors.New("missing user ID in token")
)

type JWTManager interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (string, error)
}

type jwtManager struct {
	secretKey []byte
}

func NewJWTManager(secretKey string) JWTManager {
	return &jwtManager{
		secretKey: []byte(secretKey),
	}
}

func (m *jwtManager) GenerateToken(userID string) (string, error) {
	if userID == "" {
		userIDBytes := make([]byte, 16)
		_, err := rand.Read(userIDBytes)
		if err != nil {
			return "", err
		}
		userID = hex.EncodeToString(userIDBytes)
	}

	userIDBytes, err := hex.DecodeString(userID)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, m.secretKey)
	h.Write(userIDBytes)
	sign := h.Sum(nil)

	tokenBytes := append(userIDBytes, sign...)
	return hex.EncodeToString(tokenBytes), nil
}

func (m *jwtManager) ValidateToken(tokenString string) (string, error) {
	tokenBytes, err := hex.DecodeString(tokenString)
	if err != nil {
		return "", ErrInvalidToken
	}

	if len(tokenBytes) < 48 {
		return "", ErrInvalidToken
	}

	userIDBytes := tokenBytes[:16]
	sign := tokenBytes[16:48]

	h := hmac.New(sha256.New, m.secretKey)
	h.Write(userIDBytes)
	expectedSign := h.Sum(nil)

	if !hmac.Equal(sign, expectedSign) {
		return "", ErrInvalidToken
	}

	userID := hex.EncodeToString(userIDBytes)
	if userID == "" {
		return "", ErrMissingUserID
	}

	return userID, nil
}
