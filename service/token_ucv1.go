package service

import "time"

type TokenUseCaseV1 struct {
	TokenStore
	Validity time.Duration
}

func NewTokenUseCaseV1(s TokenStore) *TokenUseCaseV1 {
	return &TokenUseCaseV1{
		TokenStore: s,
	}
}

func (uc *TokenUseCaseV1) SetValidity(value time.Duration) *TokenUseCaseV1 {
	uc.Validity = value
	return uc
}

func (uc *TokenUseCaseV1) GetTokenStore() TokenStore {
	return uc.TokenStore
}

func (uc *TokenUseCaseV1) GenerateToken(label string) (t *Token, err error) {
	t, err = uc.FindByLabel(label)
	if err == nil {
		return nil, ErrResourceAlreadyExists
	}
	t = NewToken().SetLabel(label).SetValidFor(uc.Validity)
	return uc.Save(t)
}

func (uc *TokenUseCaseV1) ValidateToken(t *Token) (err error) {
	if t == nil {
		return nil
	}
	t, err = uc.FindByValue(t.Value)
	if err != nil {
		return err
	}
	return nil
}

func (uc *TokenUseCaseV1) RemoveToken(t *Token) error {
	if t == nil {
		return nil
	}
	return uc.Remove(t.Value)
}
