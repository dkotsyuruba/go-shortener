package shortener

import "crypto/rand"

type ShortenerService interface {
	GenerateID() (string, error)
}

type RealShortenerService struct{}

func (rss *RealShortenerService) GenerateID() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}

	return string(b), nil
}

func NewRealShortenerService() ShortenerService {
	return &RealShortenerService{}
}
