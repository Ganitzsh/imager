package service_test

import (
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/ganitzsh/12fact/service"
	"github.com/stretchr/testify/require"
)

// Test cases definition

func testTokenStoreFindByValue(s service.TokenStore) func(t *testing.T) {
	return func(t *testing.T) {
		var err error

		require.NotNilf(t, s, "TokenStore is nil")
		t1 := service.NewToken().SetLabel("service_1").
			SetValidFor(24 * time.Hour).SetLabel(faker.Password())
		t2 := service.NewToken().SetLabel("service_2").
			SetValidFor(48 * time.Hour).SetLabel(faker.Password())

		t1, err = s.Save(t1)
		require.NotNil(t, t1)
		require.NoError(t, err)
		t2, err = s.Save(t2)
		require.NotNil(t, t2)
		require.NoError(t, err)

		t1Ret, err := s.FindByValue(t1.Value)
		require.NotNil(t, t1Ret)
		require.NoError(t, err)
		require.Equal(t, t1.Label, t1Ret.Label)
		require.Equal(t, t1.Value, t1Ret.Value)
		require.Equal(t, t1.ValidFor, 24*time.Hour)

		t2Ret, err := s.FindByValue(t2.Value)
		require.NotNil(t, t2Ret)
		require.NoError(t, err)
		require.Equal(t, t2.Label, t2Ret.Label)
		require.Equal(t, t2.Value, t2Ret.Value)
		require.Equal(t, t2.ValidFor, 48*time.Hour)
	}
}

func testTokenStoreFindByLabel(s service.TokenStore) func(t *testing.T) {
	return func(t *testing.T) {
		var err error

		require.NotNilf(t, s, "TokenStore is nil")
		t1 := service.NewToken().SetLabel("service_1").
			SetValidFor(24 * time.Hour).SetLabel(faker.Password())
		t2 := service.NewToken().SetLabel("service_2").
			SetValidFor(48 * time.Hour).SetLabel(faker.Password())

		t1, err = s.Save(t1)
		require.NotNil(t, t1)
		require.NoError(t, err)
		t2, err = s.Save(t2)
		require.NotNil(t, t2)
		require.NoError(t, err)

		t1Ret, err := s.FindByLabel(t1.Label)
		require.NotNil(t, t1Ret)
		require.NoError(t, err)
		require.Equal(t, t1.Label, t1Ret.Label)
		require.Equal(t, t1.Value, t1Ret.Value)
		require.Equal(t, t1.ValidFor, 24*time.Hour)

		t2Ret, err := s.FindByLabel(t2.Label)
		require.NotNil(t, t2Ret)
		require.NoError(t, err)
		require.Equal(t, t2.Label, t2Ret.Label)
		require.Equal(t, t2.Value, t2Ret.Value)
		require.Equal(t, t2.ValidFor, 48*time.Hour)
	}
}

func testTokenStoreRemove(s service.TokenStore) func(t *testing.T) {
	return func(t *testing.T) {
		var err error

		require.NotNilf(t, s, "TokenStore is nil")
		t1 := service.NewToken().SetLabel("service_1").
			SetValidFor(24 * time.Hour).SetLabel(faker.Password())
		t2 := service.NewToken().SetLabel("service_2").
			SetValidFor(48 * time.Hour).SetLabel(faker.Password())

		t1, err = s.Save(t1)
		require.NotNil(t, t1)
		require.NoError(t, err)
		t2, err = s.Save(t2)
		require.NotNil(t, t2)
		require.NoError(t, err)

		t1Ret, err := s.FindByLabel(t1.Label)
		require.NotNil(t, t1Ret)
		require.NoError(t, err)
		require.Equal(t, t1.Label, t1Ret.Label)
		require.Equal(t, t1.Value, t1Ret.Value)
		require.Equal(t, t1.ValidFor, 24*time.Hour)
		require.NoError(t, s.Remove(t1.Value))
		_, err = s.FindByValue(t1.Value)
		require.Equal(t, service.ErrResourceNotFound, err)
		_, err = s.FindByLabel(t1.Label)
		require.Equal(t, service.ErrResourceNotFound, err)

		t2Ret, err := s.FindByLabel(t2.Label)
		require.NotNil(t, t2Ret)
		require.NoError(t, err)
		require.Equal(t, t2.Label, t2Ret.Label)
		require.Equal(t, t2.Value, t2Ret.Value)
		require.Equal(t, t2.ValidFor, 48*time.Hour)
		require.NoError(t, s.Remove(t2.Value))
		_, err = s.FindByValue(t2.Value)
		require.Equal(t, service.ErrResourceNotFound, err)
		_, err = s.FindByLabel(t2.Label)
		require.Equal(t, service.ErrResourceNotFound, err)
	}
}

func testTokenStoreSave(s service.TokenStore) func(t *testing.T) {
	return func(t *testing.T) {
		var err error

		require.NotNilf(t, s, "TokenStore is nil")
		t1 := service.NewToken().SetLabel("service_1").
			SetValidFor(24 * time.Hour).SetLabel(faker.Password())
		t2 := service.NewToken().SetLabel("service_2").
			SetValidFor(48 * time.Hour).SetLabel(faker.Password())

		t1, err = s.Save(t1)
		require.NotNil(t, t1)
		require.NoError(t, err)
		t2, err = s.Save(t2)
		require.NotNil(t, t2)
		require.NoError(t, err)

		t1Ret, err := s.FindByLabel(t1.Label)
		require.NotNil(t, t1Ret)
		require.NoError(t, err)
		require.Equal(t, t1.Label, t1Ret.Label)
		require.Equal(t, t1.Value, t1Ret.Value)
		require.Equal(t, t1.ValidFor, 24*time.Hour)

		t1Ret, err = s.FindByValue(t1.Value)
		require.NotNil(t, t1Ret)
		require.NoError(t, err)
		require.Equal(t, t1.Label, t1Ret.Label)
		require.Equal(t, t1.Value, t1Ret.Value)
		require.Equal(t, t1.ValidFor, 24*time.Hour)

		t2Ret, err := s.FindByLabel(t2.Label)
		require.NotNil(t, t2Ret)
		require.NoError(t, err)
		require.Equal(t, t2.Label, t2Ret.Label)
		require.Equal(t, t2.Value, t2Ret.Value)
		require.Equal(t, t2.ValidFor, 48*time.Hour)

		t2Ret, err = s.FindByValue(t2.Value)
		require.NotNil(t, t2Ret)
		require.NoError(t, err)
		require.Equal(t, t2.Label, t2Ret.Label)
		require.Equal(t, t2.Value, t2Ret.Value)
		require.Equal(t, t2.ValidFor, 48*time.Hour)
	}
}

func testTokenUseCaseGenerateToken(
	uc service.TokenUseCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		require.NotNilf(t, uc, "TokenUseCase is nil")
	}
}

func testTokenUseCaseValidateToken(
	uc service.TokenUseCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		require.NotNilf(t, uc, "TokenUseCase is nil")
	}
}

func testTokenUseCaseRemoveToken(
	uc service.TokenUseCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		require.NotNilf(t, uc, "TokenUseCase is nil")
	}
}

// TokenStoreConsul tests

func getTokenStoreConsul() *service.TokenStoreConsul {
	return service.NewTokenStoreConsul()
}

func TestTokenStoreConsulFindByValue(t *testing.T) {
	testTokenStoreFindByValue(getTokenStoreConsul())(t)
}

func TestTokenStoreConsulFindByLabel(t *testing.T) {
	testTokenStoreFindByLabel(getTokenStoreConsul())(t)
}

func TestTokenStoreConsulRemove(t *testing.T) {
	testTokenStoreRemove(getTokenStoreConsul())(t)
}

func TestTokenStoreConsulSave(t *testing.T) {
	testTokenStoreSave(getTokenStoreConsul())(t)
}

// TokenUseCaseV1 tests

func getTokenUseCaseV1(s service.TokenStore) *service.TokenUseCaseV1 {
	return service.NewTokenUseCaseV1(s)
}

func TestTokenUseCaseV1GenerateToken(t *testing.T) {
	testTokenUseCaseGenerateToken(getTokenUseCaseV1(getTokenStoreConsul()))(t)
}

func TestTokenUseCaseV1ValidateToken(t *testing.T) {
	testTokenUseCaseValidateToken(getTokenUseCaseV1(getTokenStoreConsul()))(t)
}

func TestTokenUseCaseV1RemoveToken(t *testing.T) {
	testTokenUseCaseRemoveToken(getTokenUseCaseV1(getTokenStoreConsul()))(t)
}
