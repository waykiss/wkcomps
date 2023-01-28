package cache

import (
	"github.com/stretchr/testify/assert"
	"github.com/waykiss/wkcomps/cache/redis"
	"testing"
	"time"
)

func TestSetGetDelete(t *testing.T) {
	providers := []Provider{redis.New("127.0.0.1", "", 0)}

	for _, provider := range providers {
		t.Run(provider.GetName(), func(t *testing.T) {
			testCases := []struct {
				name, prefix, key string
				value             interface{}
				expected          interface{}
			}{
				{"string value no prefix", "", "test", "A", "A"},
			}
			for _, v := range testCases {
				t.Run(v.name, func(t *testing.T) {
					cache := New(&provider)
					cache.Prefix = v.prefix
					err := cache.Set(v.key, v.value, time.Second*10)
					assert.Nil(t, err)
					err = cache.Get(v.key, &v.value)
					assert.Nil(t, err)
					assert.Equal(t, v.expected, v.value)
				})
			}
		})
	}
}
