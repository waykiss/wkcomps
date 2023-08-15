package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/waykiss/wkcomps/cache"
	"time"
)

var redisInstance *redis.Client

type Cache struct {
	host     string
	password string
	db       int
}

func New(host, password string, dbNumber int) cache.Provider {
	return &Cache{
		host:     host,
		password: password,
		db:       dbNumber,
	}
}

func (c Cache) getRedisInstance() *redis.Client {
	if redisInstance == nil {
		// com DialTimeout de 2 segundos e MaxRetries 1, signfica que o máximo de timeout de conexão é 4 segundos
		redisInstance = redis.NewClient(&redis.Options{
			Addr:        c.host,
			Password:    c.password,
			DB:          c.db,
			DialTimeout: time.Second * 2,
			MaxRetries:  1,
		})
	}
	return redisInstance
}

// redisProcessError funcao para processar erro de conexão com o redis, notificar sobre erro e controlar quando
// foi a ultima tentativa
func redisProcessError(err error) error {
	msg := err.Error()

	// isso significa que o redis nao deu erro, mas sim nao achou a key, portanto, passe adiante, nao é um erro
	if msg == "redis: nil" {
		return nil
	}
	return nil
}

func (c Cache) Get(key string) (r []byte, err error) {
	rdb := c.getRedisInstance()
	r, err = rdb.Get(context.Background(), key).Bytes()
	if err != nil {
		err = redisProcessError(err)
	}
	return
}

func (c Cache) Set(key string, value interface{}, expiration time.Duration) (err error) {
	rdb := c.getRedisInstance()
	r := rdb.Set(context.Background(), key, value, expiration)
	if r.Err() != nil {
		redisProcessError(err)
	}
	return
}
func (c Cache) Delete(key string) error {
	return nil
}

func (Cache) GetName() string {
	return "Redis"
}
