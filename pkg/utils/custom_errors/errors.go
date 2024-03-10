package custom_errors

import "errors"

var (
	UniqueSpecialistErr = errors.New("Специалист с таким логином уже существует.")
	UniqueRatedErr      = errors.New("Вы уже оценили этот кейс.")
	NeedToAuthorizeErr  = errors.New("Необходимо заново авторизоваться")

	NoRowsCaseErr            = errors.New("Случай с таким id не найдена")
	NoRowsSpecialistLoginErr = errors.New("Пользователь с таким логином не найден")
	NoRowsSpecialistIDErr    = errors.New("Пользователь с таким id не найден")
	NoRowsCameraErr          = errors.New("Камера с таким id не найдена")

	UserUnverified = errors.New("Аккаунт пользователя не подтвержден")
	UserBadLevel   = errors.New("Аккаунт пользователя имеет неподходящий уровень")
)
