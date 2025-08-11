package util

import "golang.org/x/crypto/bcrypt"

func HashToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareTokenHash(token, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(token)) == nil
}
