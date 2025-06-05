package store

import (
	"context"
	"licenser/server/types"
	"time"

	"github.com/redis/go-redis/v9"
)

type AppStore interface {
	GetApp(ctx context.Context, name string) (*types.App, error)
	GetAppList(ctx context.Context) ([]*types.App, error)
	InsertApp(ctx context.Context, app *types.App) (*types.App, error)
}

type CachedStore struct {
	db    AppStore
	redis *redis.Client
	ttl   time.Duration
}

func NewChachedStore(db AppStore, redisClient *redis.Client, ttl time.Duration) *CachedStore {
	return &CachedStore{
		db:    db,
		redis: redisClient,
		ttl:   ttl,
	}
}
