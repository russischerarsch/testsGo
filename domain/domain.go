package domain

import "errors"

type User struct {
	Id      int
	Name    string
	Email   string
	Balance float64
}

var ErrUserAlreadyExists = errors.New("пользователь уже существует")
var ErrInvalidName = errors.New("невалидное имя")
var ErrInvalidEmail = errors.New("невалидная почта")
