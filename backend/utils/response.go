package utils

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewFailedResponse(message string) *Response {
	return &Response{
		Status:  false,
		Message: message,
		Data:    nil,
	}
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Status:  true,
		Message: "success",
		Data:    data,
	}
}
