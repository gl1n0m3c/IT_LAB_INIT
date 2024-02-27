package validators

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func CustomErrorMessage(err error) string {
	var sb strings.Builder

	// Проверяем, является ли ошибка ошибкой валидации
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				sb.WriteString(fmt.Sprintf("Поле %s является обязательным.", e.Field()))
			case "password":
				sb.WriteString("Пароль должен состоять минимум из 8 символов, заглавных и строчных букв.")
			default:
				sb.WriteString("The field " + e.Field() + " is invalid. ")
			}
		}
	} else {
		// Если ошибка не является ошибкой валидации
		sb.WriteString(err.Error())
	}

	return sb.String()
}
