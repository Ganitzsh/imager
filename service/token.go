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
}

// TokenUseCase is definition of the different actions possible regarding auth
// tokens
type TokenUseCase interface {
	GenerateToken(label string) (*Token, error)
	ValidateToken(t *Token) error
	RemoveToken(t *Token) error
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
	if err := c.Ping().Err(); err != nil {
		return nil, err
	}
	return &TokenStoreRedis{c}, nil
}

func (s *TokenStoreRedis) FindByValue(value uuid.UUID) (*Token, error) {
	return nil, nil
}

func (s *TokenStoreRedis) FindByLabel(value string) (*Token, error) {
	return nil, nil
}

func (s *TokenStoreRedis) Remove(value uuid.UUID) error {
	return nil
}

func (s *TokenStoreRedis) Save(t *Token) (*Token, error) {
	return nil, nil
}
