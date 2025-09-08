package shortener

import "crypto/rand"

type ShortenerService interface {
	GenerateID() string
}

type RealShortenerService struct{}

func (rss *RealShortenerService) GenerateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	b := make([]byte, length)
	rand.Read(b)

	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}

	return string(b)
}

func NewRealShortenerService() ShortenerService {
	return &RealShortenerService{}
}
