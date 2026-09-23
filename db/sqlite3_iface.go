package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"sync"
	"web_server_go/users"
	"web_server_go/utils"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var validTableName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// SQLite3 - хранилище пользователей на SQLite.
// Потокобезопасно, пул соединений настраивается в конструкторе.
type SQLite3 struct {
	db    *sql.DB
	table string
	once  sync.Once
}

// NewSQLite3 открывает БД по указанному пути и проверяет соединение.
// Вызывающий обязан вызвать Close, обычно через defer.
func NewSQLite3(ctx context.Context, filepath, table string) (*SQLite3, error) {
	if !validTableName.MatchString(table) {
		return nil, fmt.Errorf("NewSQLite3: invalid table name %q", table)
	}

	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, fmt.Errorf("NewSQLite3: open: %w", err)
	}

	// SQLite плохо переносит параллельную запись - ограничиваем пул.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("NewSQLite3: ping: %w", err)
	}

	return &SQLite3{db: db, table: table}, nil
}

// Close закрывает пул соединений. Идемпотентен.
func (s *SQLite3) Close() error {
	var err error
	s.once.Do(func() {
		err = s.db.Close()
	})
	return err
}

// CreateUser добавляет пользователя в БД.
// Возвращает users.ErrUserAlreadyExists при нарушении UNIQUE,
// users.ErrInvalidUserData при нарушении NOT NULL или CHECK.
func (s *SQLite3) CreateUser(ctx context.Context, client users.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(client.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("CreateUser(%s): bcrypt: %w", client.Email, err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("CreateUser(%s): begin: %w", client.Email, err)
	}
	defer tx.Rollback() // no-op после успешного Commit

	query := fmt.Sprintf(
		"INSERT INTO %s(name, last_name, father_name, email, password_hash) VALUES(?,?,?,?,?)",
		s.table,
	)
	_, err = tx.ExecContext(ctx, query,
		client.First_name, client.Last_name, client.Father_name, client.Email, string(password),
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique, sqlite3.ErrConstraintPrimaryKey:
				return fmt.Errorf("%w: %s", users.ErrUserAlreadyExists, client.Email)
			case sqlite3.ErrConstraintNotNull, sqlite3.ErrConstraintCheck:
				return fmt.Errorf("%w: %v", users.ErrInvalidUserData, err)
			}
		}
		return fmt.Errorf("CreateUser(%s): %w", client.Email, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("CreateUser(%s): commit: %w", client.Email, err)
	}
	return nil
}

// Функция получения информации об пользователе из СУБД SQLITE3.
// На вход принимает почту пользователя и возвращает структуру users.User ({Name: string , LastName: string, FatherName: string, Email: string})
func (s *SQLite3) GetUserByEmail(ctx context.Context, email string) (users.User, error) {
	var u users.User
	query := fmt.Sprintf("SELECT name, last_name, father_name, email FROM %s WHERE email = ?", s.table)
	err := s.db.QueryRowContext(ctx, query, email).Scan(&u.First_name, &u.Last_name, &u.Father_name, &u.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return u, users.ErrUserNotFoundDB
	}
	if err != nil {
		utils.Logger.Printf("Ошибка при поиске пользователя %s в базе данных (%v)\n", email, err)
		return u, fmt.Errorf("GetUserByEmail(%s): %w", email, err)
	}
	return u, nil
}

// GetPasswordHash возвращает bcrypt-хеш пароля пользователя по email.
// Если пользователь не найден — возвращает users.ErrUserNotFound.
func (s *SQLite3) GetPasswordHash(ctx context.Context, email string) ([]byte, error) {
	query := fmt.Sprintf("SELECT password_hash FROM %s WHERE email = ?", s.table)

	var storedHash string
	err := s.db.QueryRowContext(ctx, query, email).Scan(&storedHash)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %s", users.ErrUserNotFoundDB, email)
	}
	if err != nil {
		return nil, fmt.Errorf("GetPasswordHash(%s): %w", email, err)
	}

	return []byte(storedHash), nil
}
