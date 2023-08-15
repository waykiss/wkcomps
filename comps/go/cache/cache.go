package cache

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

type Provider interface {
	Set(key string, value interface{}, expirationInSeconds time.Duration) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	GetName() string
}

type Cache struct {
	provider     Provider
	lastTry      time.Time
	retrySeconds int
	Prefix       string
}

func New(provider Provider) Cache {
	return Cache{provider: provider}
}

// Get retorna um valor no redis baseado na key passada como parametro, já armazena o resultato no parametro `dest`
func (c Cache) Get(key string, dest interface{}) (err error) {
	if c.canConnect() {
		key = c.getKeyWithPrefix(key)
		start := time.Now()
		val, err := c.provider.Get(key)
		if err != nil {
			log.Errorf("error while trying to Get from provider, error: %v", err)
			return err
		}

		elapsed := time.Now().Sub(start)
		log.Info(fmt.Sprintf("Cache %s Get %s time %v", c.provider.GetName(), key, elapsed))
		if len(val) > 0 {
			err = json.Unmarshal(val, dest)
			if err != nil {
				log.Errorf("error trying to unmarshal to structure: %v", err)
			}
		}
	}
	return err
}

// Set seta um valor no redis, `expirationInSeconds` passado como zero significa que nao expira o cache
func (c Cache) Set(key string, value interface{}, expirationInSeconds time.Duration) (err error) {
	if c.canConnect() {
		key = c.getKeyWithPrefix(key)
		bytes, err := json.Marshal(value)
		if err != nil {
			log.Errorf("error trying to unmarshal to structure: %v", err)
			return err
		}
		err = c.provider.Set(key, bytes, expirationInSeconds)
		if err != nil {
			log.Errorf("error trying to unmarshal to structure: %v", err)
		}
	}
	return nil
}

// canConnect determine if it can try to connect to the provider
func (c Cache) canConnect() bool {
	if c.provider != nil {
		if c.lastTry.IsZero() {
			return true
		}
		// if it has passed more than c.retrySeconds since last attemps
		dif := time.Now().Sub(c.lastTry).Seconds()
		return int(dif) > c.retrySeconds
	}
	return false
}

// getKeyWithPrefix return the key with the prefix if was defnied
func (c Cache) getKeyWithPrefix(key string) string {
	if c.Prefix == "" {
		return key
	}
	return fmt.Sprintf("%s_%s", c.Prefix, key)
}
