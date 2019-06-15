package service_test

import (
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/ganitzsh/12fact/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	service.InitConfig()
}

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
		require.NoError(t, err)
		require.NotNil(t, t1)
		t2, err = s.Save(t2)
		require.NoError(t, err)
		require.NotNil(t, t2)

		t1Ret, err := s.FindByValue(t1.Value)
		require.NoError(t, err)
		require.NotNil(t, t1Ret)
		require.Equal(t, t1.Label, t1Ret.Label)
		require.Equal(t, t1.Value, t1Ret.Value)
		require.Equal(t, t1.ValidFor, 24*time.Hour)

		t2Ret, err := s.FindByValue(t2.Value)
		require.NoError(t, err)
		require.NotNil(t, t2Ret)
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
		require.NoError(t, err)
		require.NotNil(t, t1)
		t2, err = s.Save(t2)
		require.NoError(t, err)
		require.NotNil(t, t2)

		t1Ret, err := s.FindByLabel(t1.Label)
		require.NoError(t, err)
		require.NotNil(t, t1Ret)
		require.Equal(t, t1.Label, t1Ret.Label)
		require.Equal(t, t1.Value, t1Ret.Value)
		require.Equal(t, t1.ValidFor, 24*time.Hour)

		t2Ret, err := s.FindByLabel(t2.Label)
		require.NoError(t, err)
		require.NotNil(t, t2Ret)
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
		require.NoError(t, err)
		require.NotNil(t, t1)
		t2, err = s.Save(t2)
		require.NoError(t, err)
		require.NotNil(t, t2)

		t1Ret, err := s.FindByLabel(t1.Label)
		require.NoError(t, err)
		require.NotNil(t, t1Ret)
		require.Equal(t, t1.Label, t1Ret.Label)
		require.Equal(t, t1.Value, t1Ret.Value)
		require.Equal(t, t1.ValidFor, 24*time.Hour)
		require.NoError(t, s.Remove(t1.Value))
		_, err = s.FindByValue(t1.Value)
		require.Equal(t, service.ErrResourceNotFound, err)
		_, err = s.FindByLabel(t1.Label)
		require.Equal(t, service.ErrResourceNotFound, err)

		t2Ret, err := s.FindByLabel(t2.Label)
		require.NoError(t, err)
		require.NotNil(t, t2Ret)
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
		require.NoError(t, err)
		require.NotNil(t, t1)
		_, err = s.Save(t1)
		require.Equal(t, service.ErrResourceAlreadyExists, err)
		t2, err = s.Save(t2)
		require.NoError(t, err)
		require.NotNil(t, t2)
		_, err = s.Save(t2)
		require.Equal(t, service.ErrResourceAlreadyExists, err)

		t1Ret, err := s.FindByLabel(t1.Label)
		require.NoError(t, err)
		require.NotNil(t, t1Ret)
		require.Equal(t, t1.Label, t1Ret.Label)
		require.Equal(t, t1.Value, t1Ret.Value)
		require.Equal(t, t1.ValidFor, 24*time.Hour)

		t1Ret, err = s.FindByValue(t1.Value)
		require.NoError(t, err)
		require.NotNil(t, t1Ret)
		require.Equal(t, t1.Label, t1Ret.Label)
		require.Equal(t, t1.Value, t1Ret.Value)
		require.Equal(t, t1.ValidFor, 24*time.Hour)

		t2Ret, err := s.FindByLabel(t2.Label)
		require.NoError(t, err)
		require.NotNil(t, t2Ret)
		require.Equal(t, t2.Label, t2Ret.Label)
		require.Equal(t, t2.Value, t2Ret.Value)
		require.Equal(t, t2.ValidFor, 48*time.Hour)

		t2Ret, err = s.FindByValue(t2.Value)
		require.NoError(t, err)
		require.NotNil(t, t2Ret)
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

		label := "test." + uuid.New().String()
		t1, err := uc.GenerateToken(label)
		require.NoError(t, err)
		require.NotNil(t, t1)

		t1Dup, err := uc.GenerateToken(label)
		require.Error(t, err)
		require.Nil(t, t1Dup)

		_, err = uc.GenerateToken("")
		require.Error(t, err)
	}
}

func testTokenUseCaseValidateToken(
	uc service.TokenUseCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		require.NotNilf(t, uc, "TokenUseCase is nil")

		label := "test." + uuid.New().String()
		t1, err := uc.GenerateToken(label)
		require.NoError(t, err)
		require.NotNil(t, t1)

		require.NoError(t, uc.ValidateToken(t1))

		require.NoError(t, uc.ValidateToken(nil))

		random := service.NewToken().SetLabel("random_label")
		require.Error(t, uc.ValidateToken(random))
	}
}

func testTokenUseCaseRemoveToken(
	uc service.TokenUseCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		require.NotNilf(t, uc, "TokenUseCase is nil")

		label1 := "test." + uuid.New().String()
		t1, err := uc.GenerateToken(label1)
		require.NoError(t, err)
		require.NotNil(t, t1)

		label2 := "test." + uuid.New().String()
		t2, err := uc.GenerateToken(label2)
		require.NoError(t, err)
		require.NotNil(t, t1)

		require.NoError(t, uc.RemoveToken(t1))
		require.NoError(t, uc.RemoveToken(t2))
		require.Error(t, uc.RemoveToken(t1))
		require.Error(t, uc.RemoveToken(t2))

		require.NoError(t, uc.RemoveToken(nil))
	}
}

// TokenStoreRedis tests

func getTokenStoreRedis(t *testing.T) *service.TokenStoreRedis {
	cfg, _ := service.NewConfig()
	cfg.Store.Type = service.StoreTypeRedis
	s, err := service.GetTokenStore(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return s.(*service.TokenStoreRedis)
}

func TestTokenStoreRedisFindByValue(t *testing.T) {
	testTokenStoreFindByValue(getTokenStoreRedis(t))(t)
}

func TestTokenStoreRedisFindByLabel(t *testing.T) {
	testTokenStoreFindByLabel(getTokenStoreRedis(t))(t)
}

func TestTokenStoreRedisRemove(t *testing.T) {
	testTokenStoreRemove(getTokenStoreRedis(t))(t)
}

func TestTokenStoreRedisSave(t *testing.T) {
	testTokenStoreSave(getTokenStoreRedis(t))(t)
}

// TokenUseCaseV1 tests

func getTokenUseCaseV1(s service.TokenStore) *service.TokenUseCaseV1 {
	return service.NewTokenUseCaseV1(s)
}

func TestTokenUseCaseV1GenerateToken(t *testing.T) {
	testTokenUseCaseGenerateToken(getTokenUseCaseV1(getTokenStoreRedis(t)))(t)
}

func TestTokenUseCaseV1ValidateToken(t *testing.T) {
	testTokenUseCaseValidateToken(getTokenUseCaseV1(getTokenStoreRedis(t)))(t)
}

func TestTokenUseCaseV1RemoveToken(t *testing.T) {
	testTokenUseCaseRemoveToken(getTokenUseCaseV1(getTokenStoreRedis(t)))(t)
}
