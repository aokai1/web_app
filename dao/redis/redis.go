package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"web_app/config"
)

var rdb *redis.Client

func Init(app *config.AppConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", app.Redis.Host, app.Redis.Port),
		Password: app.Redis.Password,
		DB:       app.Redis.Database,
		PoolSize: app.Redis.PoolSize,
	})
	_, err = rdb.Ping().Result()
	return
}

func Close() {
	_ = rdb.Close()
}
