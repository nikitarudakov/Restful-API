package cache

import (
	"encoding/json"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
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

func (db *Database) SetCache(key string, value interface{}) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return db.c.Set(key, v, 1*time.Minute).Err()
}

func (db *Database) GetCache(key string, dest interface{}) error {
	v, err := db.c.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(v), dest)
}
