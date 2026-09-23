package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(plainPassword []byte, hashedPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(plainPassword))
	Logger.Printf("plainPassword is %q , hashedPassword is %q\n", plainPassword, hashedPassword)
	if err != nil {
		Logger.Printf("Неверно введен пароль (%v)\n", err)
		return false
	}
	return true
}
