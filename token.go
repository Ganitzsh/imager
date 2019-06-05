package main

import (
	"time"

	"github.com/hashicorp/consul/api"
)

type Token struct {
	CreatedAt time.Time
	ValidFor  time.Duration
	Label     string
	Value     string
}

func NewToken() *Token {
	return &Token{}
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

func (t *Token) SetValue(value string) *Token {
	t.Value = value
	return t
}

// TokenStore is an interface defining the behaviour of the store that will be
// responsible for manipulating the authentication tokens from the storage
type TokenStore interface {
	FindByValue(value string) (*Token, error)
	FindByLabel(label string) (*Token, error)
	Remove(value string) error
	Save(t *Token) (*Token, error)
}

// TokenUseCase is definition of the different actions possible regarding auth
// tokens
type TokenUseCase interface {
	GenerateToken(label string) (*Token, error)
	ValidateToken(value string) error
	RemoveToken(value string) error
}

// TokenStoreConsul is an implemetation of the TokenStore backed by the
// key-value store Consul
type TokenStoreConsul struct {
	*api.KV
}

func (s *TokenStoreConsul) FindByValue(value string) (*Token, error) {
	return nil, nil
}

func (s *TokenStoreConsul) FindByLabel(value string) (*Token, error) {
	return nil, nil
}

func (s *TokenStoreConsul) Remove(value string) error {
	return nil
}

func (s *TokenStoreConsul) Save(t *Token) (*Token, error) {
	return nil, nil
}
