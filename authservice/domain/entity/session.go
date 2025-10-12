package entity

import (
	"errors"
	"strings"
	"time"
)

// ===== Value Objects =====

type UserID uint64

func (id UserID) Valid() bool { return id != 0 }

type JTI string

func (jti JTI) Valid() bool { return strings.TrimSpace(string(jti)) != "" }

// ===== Aggregate Root =====

type Session struct {
	// поля не экспортируем: управление только методами
	ID          uint64
	userID      UserID
	jti         JTI
	issuedAt    time.Time
	expiresAt   time.Time
	rotatedFrom *string
	userAgent   *string
	ip          *string
	revokedAt   time.Time
	reusedAt    time.Time
}

var (
	ErrInvalidSessionID      = errors.New("session: invalid session ID")
	ErrInvalidSUserID        = errors.New("session: invalid user ID")
	ErrInvalidJTI            = errors.New("session: invalid JTI")
	ErrSessionExpired        = errors.New("session: session expired")
	ErrSessionRevoked        = errors.New("session: session revoked")
	ErrSessionAlreadyRevoked = errors.New("session: session already revoked")
)

// конструктор

func NewSession(id uint64, userID UserID, jti string, issuedAt, expiresAt time.Time, opts ...Option) (*Session, error) {
	if id == 0 {
		return nil, ErrInvalidSessionID
	}
	if !userID.Valid() {
		return nil, ErrInvalidSUserID
	}
	if !JTI(jti).Valid() {
		return nil, ErrInvalidJTI
	}

	session := &Session{
		ID:          id,
		userID:      userID,
		jti:         JTI(jti),
		issuedAt:    issuedAt,
		expiresAt:   expiresAt,
		rotatedFrom: nil,
		userAgent:   nil,
		ip:          nil,
		revokedAt:   time.Time{},
		reusedAt:    time.Time{},
	}

	// обработка опций
	for _, opt := range opts {
		if err := opt(session); err != nil {
			return nil, err
		}
	}

	return session, nil
}

// функциональные опции для необязательных полей

type Option func(*Session) error

func WithUserAgent(userAgent string) Option {
	return func(s *Session) error {
		s.userAgent = &userAgent
		return nil
	}
}

func WithIP(ip string) Option {
	return func(s *Session) error {
		s.ip = &ip
		return nil
	}
}

func WithRotatedFrom(rotatedFrom string) Option {
	return func(s *Session) error {
		s.rotatedFrom = &rotatedFrom
		return nil
	}
}

// ===== Геттеры (только чтение) =====

func (s *Session) SessionID() uint64    { return s.ID }
func (s *Session) UserID() UserID       { return s.userID }
func (s *Session) JTI() JTI             { return s.jti }
func (s *Session) IssuedAt() time.Time  { return s.issuedAt }
func (s *Session) ExpiresAt() time.Time { return s.expiresAt }
func (s *Session) RotatedFrom() *string { return s.rotatedFrom }
func (s *Session) UserAgent() *string   { return s.userAgent }
func (s *Session) IP() *string          { return s.ip }
func (s *Session) RevokedAt() time.Time { return s.revokedAt }
func (s *Session) ReusedAt() time.Time  { return s.reusedAt }

// ===== Командные методы (мутация через поведение) =====

func (s *Session) Revoke() error {
	if !s.revokedAt.IsZero() {
		return ErrSessionAlreadyRevoked
	}
	if s.IsExpired() {
		return ErrSessionExpired
	}

	s.revokedAt = time.Now().UTC()
	return nil
}

func (s *Session) Reuse() error {
	if s.IsRevoked() {
		return ErrSessionRevoked
	}
	if s.IsExpired() {
		return ErrSessionExpired
	}

	s.reusedAt = time.Now().UTC()
	return nil
}

func (s *Session) Rotate(newJTI string) error {
	if s.IsRevoked() {
		return ErrSessionRevoked
	}
	if s.IsExpired() {
		return ErrSessionExpired
	}

	// создаем новую сессию на основе старой
	s.rotatedFrom = (*string)(&s.jti)
	s.jti = JTI(newJTI)
	s.issuedAt = time.Now().UTC()
	return nil
}

func (s *Session) IsExpired() bool {
	return time.Now().UTC().After(s.expiresAt)
}

func (s *Session) IsRevoked() bool {
	return !s.revokedAt.IsZero()
}

func (s *Session) touch() *Session {
	s.issuedAt = time.Now().UTC()
	s.expiresAt = s.issuedAt.Add(24 * time.Hour)
	return s
}
