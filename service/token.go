package service

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type Token struct {
	CreatedAt time.Time
	ValidFor  time.Duration
	Label     string
	Value     uuid.UUID
}

func NewToken() *Token {
	return &Token{
		CreatedAt: time.Now(),
		Value:     uuid.New(),
	}
}

func (t *Token) SetCreatedAt(value time.Time) *Token {
	t.CreatedAt = value
	return t
}

func (t *Token) SetValidFor(value time.Duration) *Token {
	t.ValidFor = value
	return t
}

func (t *Token) SetLabel(value string) *Token {
	t.Label = value
	return t
}

func (t *Token) SetValue(value uuid.UUID) *Token {
	t.Value = value
	return t
}

// TokenStore is an interface defining the behaviour of the store that will be
// responsible for manipulating the authentication tokens from the storage
type TokenStore interface {
	FindByValue(value uuid.UUID) (*Token, error)
	FindByLabel(label string) (*Token, error)
	Remove(value uuid.UUID) error
	Save(t *Token) (*Token, error)
	Up() bool
	RootToken() *Token
}

// TokenUseCase is definition of the different actions possible regarding auth
// tokens
type TokenUseCase interface {
	GenerateToken(label string) (*Token, error)
	ValidateToken(t *Token) error
	RemoveToken(t *Token) error
	GetTokenStore() TokenStore
}

// TokenStoreRedis is an implemetation of the TokenStore backed by the
// key-value store Redis
type TokenStoreRedis struct {
	*redis.Client
}

func NewTokenStoreRedis(c *redis.Client) (*TokenStoreRedis, error) {
	if c == nil {
		return nil, errors.New("Invalid redis client")
	}
	return &TokenStoreRedis{c}, nil
}

func (TokenStoreRedis) RootToken() *Token {
	return NewToken().SetLabel("transformer.root.token")
}

func (s *TokenStoreRedis) FindByValue(value uuid.UUID) (*Token, error) {
	cmd := s.Get(value.String())
	if err := cmd.Err(); err != nil {
		return nil, ErrResourceNotFound
	}
	label, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	return &Token{
		Label: label,
		Value: value,
	}, nil
}

func (s *TokenStoreRedis) FindByLabel(label string) (*Token, error) {
	cmd := s.Get(label)
	if e := cmd.Err(); e != nil {
		return nil, ErrResourceNotFound
	}
	vStr, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	v, err := uuid.Parse(vStr)
	if err != nil {
		return nil, err
	}
	return &Token{
		Label: label,
		Value: v,
	}, nil
}

func (s *TokenStoreRedis) Remove(value uuid.UUID) error {
	cmd := s.Get(value.String())
	if err := cmd.Err(); err != nil {
		return err
	}
	label, err := cmd.Result()
	if err != nil {
		return err
	}
	if err := s.Del(value.String()).Err(); err != nil {
		return err
	}
	if err := s.Del(label).Err(); err != nil {
		return err
	}
	return nil
}

func (s *TokenStoreRedis) Save(t *Token) (*Token, error) {
	if t == nil {
		return nil, ErrInternalError
	}
	if t.Label == "" {
		return nil, ErrInvalidInput
	}
	if _, err := s.FindByLabel(t.Label); err == nil {
		return nil, ErrResourceAlreadyExists
	}
	if err := s.Set(t.Label, t.Value.String(), t.ValidFor).Err(); err != nil {
		return nil, err
	}
	if err := s.Set(t.Value.String(), t.Label, t.ValidFor).Err(); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TokenStoreRedis) Up() bool {
	if s.Client == nil {
		return false
	}
	if err := s.Ping().Err(); err != nil {
		return false
	}
	return true
}
