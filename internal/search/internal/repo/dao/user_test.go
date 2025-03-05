package dao

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestElasticUserDAO_SearchByUsername(t *testing.T) {
	client := initEs(t)
	err := InitEs(client)
	require.NoError(t, err)
	dao := NewElasticUserDAO(client)
	err = dao.InputUser(context.Background(), User{
		Id:       100,
		Account:  "HelloKitty1111",
		Username: "adsadasid",
	})
	require.NoError(t, err)

	res, err := dao.SearchByUsername(context.Background(), []string{
		"hello",
	})
	require.NoError(t, err)
	t.Log(res)
}

func TestElasticUserDAO_InputUser(t *testing.T) {
	client := initEs(t)
	err := InitEs(client)
	require.NoError(t, err)
	dao := NewElasticUserDAO(client)
	err = dao.InputUser(context.Background(), User{
		Id:       101,
		Account:  "HelloKitty",
		Username: "adsadasid",
	})
	require.NoError(t, err)
}
