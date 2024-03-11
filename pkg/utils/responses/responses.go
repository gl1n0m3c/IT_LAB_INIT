package responses

const (
	Response400 = "Bad request: %s"
	Response500 = "Internal server error"

	ResponseSuccessCreate = "Объект %s был успешно создан с id: %d"
	ResponseSuccessGet    = "Объект(ы) %s с был(и) успешно получен(ы)"
	ResponseSuccessUpdate = "Объект %s с был успешно обновлен"

	ResponseNoByteStringProvided = "Байтовая строка отсутствует"
	ResponseNoPhotoProvided      = "Фото отсуствует"

	ResponseBadFileSize   = "Ваш файл слишком большой"
	ResponseBadPhotoFile  = "Вы загрузили не фото"
	ResponseBadByteString = "Байтовая строка некорректна"
	ResponseBadTime       = "Переданное время некорректно"

	ResponseBadQuery = "Параметры запроса указаны некорректно"

	ResponseSuccessDelete = "Объект %s с id %d был успешно удален"
)

type CreationIntResponse struct {
	ID int `json:"id"`
}

type CreationStringResponse struct {
	ID string `json:"id"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type MessageDataResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type JWTRefresh struct {
	JWT          string `json:"JWT"`
	RefreshToken string `json:"RefreshToken"`
}

func NewMessageResponse(message string) MessageResponse {
	return MessageResponse{Message: message}
}

func NewJWTRefreshResponse(JWT, RefreshToken string) JWTRefresh {
	return JWTRefresh{
		JWT:          JWT,
		RefreshToken: RefreshToken,
	}
}
