package users

import "errors"

var (
	ErrUserNotFoundDB    = errors.New("Пользователь не найден в базе данных")
	ErrPasswordInvalid   = errors.New("Пароль пользователя неверный")
	ErrUserAlreadyExists = errors.New("Пользователь уже существует в базе данных")
	ErrInvalidUserData   = errors.New("Неверные входящие данные для базы данных (insertion)")
)
