package gorse

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/recommend/internal/domain"
	client "github.com/gorse-io/gorse-go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSyncService(t *testing.T) {
	svc := NewSyncService(client.NewGorseClient("http://localhost:8088", ""))
	err := svc.InputUser(context.Background(), domain.User{
		Id: 1,
	})
	require.NoError(t, err)
	err = svc.InputInk(context.Background(), domain.Ink{
		Id:        1,
		Tags:      []string{"web3", "eth", "btc", "zk"},
		Title:     "test",
		CreatedAt: time.Now(),
	})
	require.NoError(t, err)

	err = svc.InputInk(context.Background(), domain.Ink{
		Id:        2,
		Tags:      []string{"btc", "polka", "dot"},
		Title:     "buy the deep",
		CreatedAt: time.Now(),
	})
}
