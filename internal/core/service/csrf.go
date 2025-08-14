package service

import (
	"crypto/rand"
	"encoding/hex"
)

type CSRFService struct{}

func NewCSRFService() *CSRFService {
	return &CSRFService{}
}

func (s *CSRFService) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *CSRFService) VerifyToken(token, expectedToken string) bool {
	return token == expectedToken
}
