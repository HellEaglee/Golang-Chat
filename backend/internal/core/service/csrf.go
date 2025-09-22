package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"time"
)

type CSRFService struct{}

func NewCSRFService() *CSRFService {
	return &CSRFService{}
}

func (s *CSRFService) GenerateToken() (string, error) {
	type csrfData struct {
		Token  string `json:"token"`
		Expire int64  `json:"expo"`
	}
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	data := csrfData{
		Token:  hex.EncodeToString(bytes),
		Expire: time.Now().Add(time.Minute).Unix(),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(jsonData), nil
}

func (s *CSRFService) VerifyToken(token, expectedToken string) bool {
	tokenBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	expectedBytes, err := base64.StdEncoding.DecodeString(expectedToken)
	if err != nil {
		return false
	}

	var tokenData, expectedData map[string]interface{}
	if err := json.Unmarshal(tokenBytes, &tokenData); err != nil {
		return false
	}
	if err := json.Unmarshal(expectedBytes, &expectedData); err != nil {
		return false
	}

	if exp, ok := tokenData["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return false
		}
	}

	tokenStr, _ := tokenData["token"].(string)
	expectedStr, _ := expectedData["token"].(string)

	tokenBytesFinal, _ := hex.DecodeString(tokenStr)
	expectedBytesFinal, _ := hex.DecodeString(expectedStr)

	return subtle.ConstantTimeCompare(tokenBytesFinal, expectedBytesFinal) == 1
}
