package api

type ErrResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewErrResponse(err error, code int) *ErrResponse {
	return &ErrResponse{
		Error{
			Message: err.Error(),
			Code:    code,
		},
	}
}
