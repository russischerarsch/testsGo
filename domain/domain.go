package domain

import "errors"

type User struct {
	Id       int
	Name     string
	Email    string
	Balance  float64
	Password string
	Age      int
}

var ErrUserAlreadyExists = errors.New("пользователь уже существует")
var ErrInvalidName = errors.New("невалидное имя")
var ErrInvalidEmail = errors.New("невалидная почта")
var ErrNoUpperLetter = errors.New("необходима буква в верхнем регистре")
var ErrShortPassword = errors.New("Слишком короткий пароль")
var ErrAgeForbidden = errors.New("пользователь должен быть старше 18 лет")
var ErrNoSpecialChar = errors.New("нет специального знака")
