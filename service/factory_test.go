package service_test

import (
	"testing"

	"github.com/ganitzsh/12fact/service"
	"github.com/stretchr/testify/require"
)

func TestGetTokenStore(t *testing.T) {
	var s service.TokenStore
	var err error

	s, err = service.GetTokenStore(nil)
	require.Error(t, err)
	require.Nil(t, s)

	s, err = service.GetTokenStore(&service.Config{
		Store: &service.StoreConfig{
			Type: service.StoreType("some_rand String"),
		},
	})
	require.Error(t, err)
	require.Nil(t, s)
}
