package custom_errors

import "errors"

var (
	UniqueSpecialistErr = errors.New("Ошибка при регистрации специалиста: специалист с таким логином уже существует.")
	NeedToAuthorizeErr  = errors.New("Необходимо заново авторизоваться")

	NoRowsLoginErr  = errors.New("Пользователь с таким логином не найден")
	NoRowsCameraErr = errors.New("Камера с таким id не найдена")
)
