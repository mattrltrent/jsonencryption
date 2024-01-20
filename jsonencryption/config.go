package jsonencryption

import "errors"

var (
	ErrInvalidKey = errors.New("invalid key, it must be 16 bytes long")
	ErrNoSecret   = errors.New("secret object is nil, you must call SetKey() first")
)
