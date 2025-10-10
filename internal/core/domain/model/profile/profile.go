package profile

import (
	"ads/internal/pkg/ddd"
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

// ===== Value Objects =====

type UserID uint64

func (id UserID) Valid() bool { return id != 0 }

// примитивная валидация. при необходимости заверните libphonenumber.
type Phone string

func NewPhone(raw string) (Phone, error) {
	p := Phone(strings.TrimSpace(raw))
	if p == "" {
		return "", errors.New("phone: empty")
	}
	// мини-инвариант: только цифры, +, -, пробел, скобки
	for _, r := range p {
		if !(r == '+' || r == '-' || r == ' ' || r == '(' || r == ')' || (r >= '0' && r <= '9')) {
			return "", errors.New("phone: invalid characters")
		}
	}
	return p, nil
}

type Subscription string

func (s Subscription) Valid() bool { return strings.TrimSpace(string(s)) != "" }

// ===== Aggregate Root =====

type Profile struct {
	*ddd.BaseAggregate[uuid.UUID]

	// поля не экспортируем: управление только методами
	userID               UserID
	name                 string
	phone                *Phone
	photoID              *string
	notificationsEnabled bool
	subscriptions        []Subscription
	updatedAt            time.Time
	banned               bool
}

// инварианты/ошибки
var (
	ErrInvalidUserID      = errors.New("profile: invalid user id")
	ErrEmptyName          = errors.New("profile: empty name")
	ErrAlreadyBanned      = errors.New("profile: already banned")
	ErrNotBanned          = errors.New("profile: not banned")
	ErrSubscriptionExists = errors.New("profile: subscription already exists")
	ErrNoSuchSubscription = errors.New("profile: no such subscription")
)

// конструктор

func New(id UserID, name, rawPhone string, opts ...Option) (*Profile, error) {
	if !id.Valid() {
		return nil, ErrInvalidUserID
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyName
	}

	phone, err := NewPhone(rawPhone)
	if err != nil {
		return nil, err
	}

	p := &Profile{
		BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](uuid.New()),

		userID:               id,
		name:                 name,
		phone:                &phone,
		notificationsEnabled: true,
		subscriptions:        make([]Subscription, 0),
		updatedAt:            time.Now().UTC(),
		banned:               false,
	}
	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, err
		}
	}
	return p.touch(), nil
}

// функциональные опции для необязательных полей

type Option func(*Profile) error

func WithPhone(ph Phone) Option {
	return func(p *Profile) error { p.phone = &ph; return nil }
}
func WithPhotoID(photo string) Option {
	return func(p *Profile) error {
		ph := strings.TrimSpace(photo)
		if ph == "" {
			return errors.New("photoID: empty")
		}
		p.photoID = &ph
		return nil
	}
}
func WithSubscriptions(subs ...Subscription) Option {
	return func(p *Profile) error {
		for _, s := range subs {
			if !s.Valid() {
				return errors.New("invalid subscription")
			}
		}
		p.subscriptions = dedup(append(p.subscriptions, subs...))
		return nil
	}
}
func WithNotificationsEnabled(enabled bool) Option {
	return func(p *Profile) error { p.notificationsEnabled = enabled; return nil }
}

// ===== Геттеры (только чтение) =====

func (p *Profile) UserID() UserID             { return p.userID }
func (p *Profile) Name() string               { return p.name }
func (p *Profile) Phone() *Phone              { return p.phone }
func (p *Profile) PhotoID() *string           { return p.photoID }
func (p *Profile) NotificationsEnabled() bool { return p.notificationsEnabled }
func (p *Profile) Subscriptions() []Subscription {
	return append([]Subscription(nil), p.subscriptions...)
}
func (p *Profile) UpdatedAt() time.Time { return p.updatedAt }
func (p *Profile) Banned() bool         { return p.banned }

// ===== Командные методы (мутация через поведение) =====

func (p *Profile) Rename(newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return ErrEmptyName
	}
	p.name = newName
	return p.touchErr()
}

func (p *Profile) SetPhone(ph Phone) error { p.phone = &ph; return p.touchErr() }
func (p *Profile) ClearPhone() error       { p.phone = nil; return p.touchErr() }

func (p *Profile) SetPhoto(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("photoID: empty")
	}
	p.photoID = &id
	return p.touchErr()
}
func (p *Profile) ClearPhoto() error { p.photoID = nil; return p.touchErr() }

func (p *Profile) EnableNotifications() error  { p.notificationsEnabled = true; return p.touchErr() }
func (p *Profile) DisableNotifications() error { p.notificationsEnabled = false; return p.touchErr() }

func (p *Profile) AddSubscription(s Subscription) error {
	if !s.Valid() {
		return errors.New("invalid subscription")
	}
	for _, x := range p.subscriptions {
		if x == s {
			return ErrSubscriptionExists
		}
	}
	p.subscriptions = append(p.subscriptions, s)
	return p.touchErr()
}

func (p *Profile) RemoveSubscription(s Subscription) error {
	i := -1
	for idx, x := range p.subscriptions {
		if x == s {
			i = idx
			break
		}
	}
	if i == -1 {
		return ErrNoSuchSubscription
	}
	p.subscriptions = append(p.subscriptions[:i], p.subscriptions[i+1:]...)
	return p.touchErr()
}

func (p *Profile) Ban() error {
	if p.banned {
		return ErrAlreadyBanned
	}
	p.banned = true
	return p.touchErr()
}

func (p *Profile) Unban() error {
	if !p.banned {
		return ErrNotBanned
	}
	p.banned = false
	return p.touchErr()
}

// доменная «служебка»
func (p *Profile) touch() *Profile { p.updatedAt = time.Now().UTC(); return p }
func (p *Profile) touchErr() error { p.updatedAt = time.Now().UTC(); return nil }
func dedup(in []Subscription) []Subscription {
	seen := make(map[Subscription]struct{}, len(in))
	out := make([]Subscription, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
