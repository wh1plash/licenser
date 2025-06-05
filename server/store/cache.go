package store

import (
	"context"
	"encoding/json"
	"fmt"
	"licenser/server/types"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(host string, port string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
		DB:   0,
	})
}

func (c *CachedStore) GetApp(ctx context.Context, name string) (*types.App, error) {
	key := "app:" + name

	val, err := c.redis.Get(ctx, key).Result()
	if err == nil {
		var app types.App
		if err := json.Unmarshal([]byte(val), &app); err == nil {
			fmt.Println("get from Redis")
			return &app, nil
		}
	}

	fmt.Println("Redis miss")
	app, err := c.db.GetApp(ctx, name)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(app)
	c.redis.Set(ctx, key, data, c.ttl)

	return app, nil
}

func (c *CachedStore) GetAppList(ctx context.Context) ([]*types.App, error) {
	const key = "app:list"

	val, err := c.redis.Get(ctx, key).Result()
	if err == nil {
		var apps []*types.App
		if err := json.Unmarshal([]byte(val), &apps); err == nil {
			fmt.Println("get from Redis")
			return apps, nil
		}
	}
	fmt.Println("Redis miss")
	apps, err := c.db.GetAppList(ctx)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(apps)
	c.redis.Set(ctx, key, data, c.ttl)

	return apps, nil
}

func (c *CachedStore) InsertApp(ctx context.Context, app *types.App) (*types.App, error) {
	inserted, err := c.db.InsertApp(ctx, app)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	c.redis.Del(ctx, "app:"+inserted.Name)
	c.redis.Del(ctx, "app:list")

	return inserted, nil
}
