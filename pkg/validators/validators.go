package validators

import (
	"github.com/go-playground/validator/v10"
	"mime/multipart"
	"path/filepath"
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

func ValidateFileTypeExtension(file *multipart.FileHeader) bool {
	// Проверка на допустимый тип `Content-Type`
	allowedTypes := map[string]bool{
		"image/jpeg":    true,
		"image/png":     true,
		"image/svg+xml": true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		return false
	}

	extension := filepath.Ext(file.Filename)
	// Проверка на допустимое расширение файла
	allowedExtensions := map[string]bool{
		".jpeg": true,
		".jpg":  true,
		".png":  true,
		".svg":  true,
	}
	if !allowedExtensions[extension] {
		return false
	}

	return true
}
