package redis

import (
	storageT "github.com/fd/simplex/data/storage/testing"
	"github.com/simonz05/godis/redis"
	"testing"
)

func TestRedis(t *testing.T) {
	storageT.ValidateDriver(t, &S{
		Client: redis.New("", 1, ""),
	})
}
