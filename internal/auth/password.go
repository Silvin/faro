package auth

import "golang.org/x/crypto/bcrypt"

// bcryptCost equilibra seguridad y rendimiento (ver ADR-002).
const bcryptCost = 12

func hashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func checkPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
