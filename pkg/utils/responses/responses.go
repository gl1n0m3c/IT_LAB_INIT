package responses

const (
	Response201 = "Object %s was successfully created with id %d"

	Response400 = "Goddamn, that was really bad request: %s"

	Response500 = "Internal server error"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type MessageDataResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewMessageResponse(message string) MessageResponse {
	return MessageResponse{Message: message}
}

func NewMessageDataResponse(message string, data any) MessageDataResponse {
	return MessageDataResponse{Message: message, Data: data}
}
