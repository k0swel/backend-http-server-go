package users

import "context"

// import "web_server_go/internal/users"

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (User, error) // получить пользователя по email
	CreateUser(ctx context.Context, account User) error             // зарегистрировать пользователя
}
