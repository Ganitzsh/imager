package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
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

// TokenStoreConsul is an implemetation of the TokenStore backed by the
// key-value store Consul
type TokenStoreConsul struct {
	*api.KV
}

func NewTokenStoreConsul() *TokenStoreConsul {
	return &TokenStoreConsul{}
}

func (s *TokenStoreConsul) FindByValue(value uuid.UUID) (*Token, error) {
	return nil, nil
}

func (s *TokenStoreConsul) FindByLabel(value string) (*Token, error) {
	return nil, nil
}

func (s *TokenStoreConsul) Remove(value uuid.UUID) error {
	return nil
}

func (s *TokenStoreConsul) Save(t *Token) (*Token, error) {
	return nil, nil
}
