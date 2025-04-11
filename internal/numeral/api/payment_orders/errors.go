package payment_orders

import "errors"

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrInvalidRequest      = errors.New("invalid request")
	ErrRecordAlreadyExists = errors.New("record already exists")
)
