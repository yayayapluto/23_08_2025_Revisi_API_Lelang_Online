package domain

import "errors"

var (
	ErrInvalidSortByColumn = errors.New("invalid sort by column")
	ErrInvalidSortDir      = errors.New("invalid sort direction")
	ErrInvalidPayload      = errors.New("invalid payload")
)
