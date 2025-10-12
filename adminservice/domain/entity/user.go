package entity

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidEmail     = errors.New("user: invalid email")
	ErrInvalidPassword  = errors.New("user: invalid password")
	ErrUserBanned       = errors.New("user: user is banned")
	ErrUserNotBanned    = errors.New("user: user is not banned")
	ErrEmailEmpty       = errors.New("user: email is empty")
	ErrPasswordTooShort = errors.New("user: password is too short")
)

// ===== User Entity =====

type User struct {
	ID        uint64
	Email     string
	Password  string
	Role      string
	Banned    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Конструктор для создания нового пользователя
func NewUser(id uint64, email, password, role string) (*User, error) {
	// Валидация Email
	if email == "" || !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}

	// Валидация пароля
	if len(password) < 8 {
		return nil, ErrPasswordTooShort
	}

	user := &User{
		ID:        id,
		Email:     email,
		Password:  password,
		Role:      role,
		Banned:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}

// ===== Методы для изменения состояния =====

// Изменение роли пользователя
func (u *User) ChangeRole(newRole string) {
	u.Role = newRole
	u.UpdatedAt = time.Now()
}

// Бан пользователя
func (u *User) Ban() error {
	if u.Banned {
		return ErrUserBanned
	}
	u.Banned = true
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) Unban() error {
	if !u.Banned {
		return ErrUserNotBanned
	}
	u.Banned = false
	u.UpdatedAt = time.Now()
	return nil
}

// Проверка, забанен ли пользователь

func (u *User) IsBanned() bool {
	return u.Banned
}

// ===== Геттеры (только чтение) =====

func (u *User) GetID() uint64           { return u.ID }
func (u *User) GetEmail() string        { return u.Email }
func (u *User) GetPassword() string     { return u.Password }
func (u *User) GetRole() string         { return u.Role }
func (u *User) GetCreatedAt() time.Time { return u.CreatedAt }
func (u *User) GetUpdatedAt() time.Time { return u.UpdatedAt }
