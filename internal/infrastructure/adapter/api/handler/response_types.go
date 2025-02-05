package handler

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func ErrorResponse(err string) Response {
	return Response{
		Error: err,
	}
}

func SuccessResponse(data interface{}) Response {
	return Response{
		Data: data,
	}
}

func MessageResponse(message string) Response {
	return Response{
		Message: message,
	}
}
