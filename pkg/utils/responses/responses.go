package responses

const (
	Response201 = "Объект %s был успешно создан с id: %d"

	Response400 = "Bad request: %s"

	Response500 = "Internal server error"

	ResponseBadFileSize   = "Ваш файл слишком большой"
	ResponseBadFileType   = "Вы загрузили не фото"
	ResponseSuccessDelete = "Объект %s с id %d был успешно удален"
)

type CreationResponse struct {
	ID int `json:"id"`
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
