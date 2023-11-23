package cache

import (
	"context"
	"encoding/json"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/tools/timeparse"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type Database struct {
	c *redis.Client
}

func NewCacheDatabase(cfg *config.Cache) (*Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Server + ":" + cfg.Port,
		Password: "",
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	log.Info().Msg("Redis is running at http://" + cfg.Server + ":" + cfg.Port)

	return &Database{c: client}, nil
}

func Middleware(db *Database, dest interface{}, cacheCfg *config.Cache) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			expiration, err := timeparse.ParseExpirationTime(cacheCfg.ExpirationQuan,
				cacheCfg.ExpirationUnit)

			if err != nil {
				log.Error().Err(err).
					Msg("Cache middleware can't be set up because of invalid expiration time (default: 1 min.")
				expiration = 1 * time.Minute
			}

			if ctx.Request().Method != http.MethodGet {
				if err := next(ctx); err != nil {
					return err
				}
			}

			// Create a unique cache key based on method and URI
			cacheKey := ctx.Request().Method + ":" + ctx.Request().RequestURI

			if err = db.GetCache(cacheKey, dest); err == nil {
				return ctx.JSON(200, dest)
			} else {
				log.Warn().Err(err).Msg("cache was not retrieved")
			}

			if err = next(ctx); err != nil {
				return err
			}

			value := ctx.Get(cacheKey)
			if err = db.SetCache(cacheKey, value, expiration); err != nil {
				log.Warn().Str("service", "caching").Err(err).Send()
			}

			return nil
		}
	}
}

func (db *Database) SetCache(key string, value interface{}, expiration time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return db.c.Set(context.Background(), key, v, expiration).Err()
}

func (db *Database) GetCache(key string, dest interface{}) error {
	v, err := db.c.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(v), dest)
}
