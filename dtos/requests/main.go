package requests

type Response struct {
	Status  bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (r *Response) Success(message string, data any) *Response {
	return &Response{
		Status:  true,
		Message: message,
		Data:    data,
	}
}

func (r *Response) Error(message string) *Response {
	return &Response{
		Status:  false,
		Message: message,
		Data:    nil,
	}
}

func (r *Response) ErrorWithData(message string, data any) *Response {
	return &Response{
		Status:  false,
		Message: message,
		Data:    data,
	}
}
