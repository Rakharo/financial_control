package utils

import (
	"errors"
	"unicode"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("A senha deve ter pelo menos 8 caracteres")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("A senha deve conter pelo menos uma letra maiúscula")
	}
	if !hasLower {
		return errors.New("A senha deve conter pelo menos uma letra minúscula")
	}
	if !hasNumber {
		return errors.New("A senha deve conter pelo menos um número")
	}
	if !hasSpecial {
		return errors.New("A senha deve conter pelo menos um caractere especial")
	}

	return nil
}
