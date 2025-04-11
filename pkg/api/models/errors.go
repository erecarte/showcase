package models

type ApiError struct {
	Message string `json:"error"`
}

func (e ApiError) Error() string {
	return e.Message
}
