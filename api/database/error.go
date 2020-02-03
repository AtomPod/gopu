package database

import "errors"

var (
	//ErrNotFound  not found
	ErrNotFound = errors.New("not found")
	//ErrUsernameAlreadyExists username already exists
	ErrUsernameAlreadyExists = errors.New("username already exists")
	//ErrEmailAlreadyExists email already exists
	ErrEmailAlreadyExists = errors.New("email already exists")
)
