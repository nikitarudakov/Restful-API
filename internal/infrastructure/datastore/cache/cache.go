package cache

import (
	"encoding/json"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/tools/timeparse"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
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

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	log.Info().Msg("Redis is running at http://" + cfg.Server + ":" + cfg.Port)

	return &Database{c: client}, nil
}

func Middleware(db *Database, key string, dest interface{}, cacheCfg *config.Cache) echo.MiddlewareFunc {
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
				return nil
			}

			if err := db.GetCache(key, dest); err == nil {
				return ctx.JSON(200, dest)
			} else {
				log.Warn().Err(err).Msg("cache was not retrieved")
			}

			if err := next(ctx); err != nil {
				return err
			}

			value := ctx.Get(key)
			if err := db.SetCache("rating", value, expiration); err != nil {
				log.Warn().Str("service", "rating caching").Err(err).Send()
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
	return db.c.Set(key, v, expiration).Err()
}

func (db *Database) GetCache(key string, dest interface{}) error {
	v, err := db.c.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(v), dest)
}
