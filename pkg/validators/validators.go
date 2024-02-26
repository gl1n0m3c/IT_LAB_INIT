package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString
	if !hasNumber(password) {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString
	if !hasUpper(password) {
		return false
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString
	if !hasLower(password) {
		return false
	}

	return true
}
