package model

type (
	ApiRequest struct {
	}

	ApiResponse struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Error   string      `json:"error,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}
)
