package slidingLog_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/distuurbia/rateLimiters/slidingLog"
	"github.com/go-redis/redis"
	"github.com/ory/dockertest"
)

var sl *slidingLog.SlidingLog

func SetupRedis() (*redis.Client, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %w", err)
	}
	resource, err := pool.Run("redis", "latest", []string{})
	if err != nil {
		return nil, nil, fmt.Errorf("could not run the pool: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
		DB:   0,
	})
	cleanup := func() {
		client.Close()
		pool.Purge(resource)
	}
	return client, cleanup, nil
}

func TestMain(m *testing.M) {
	client, cleanupRds, err := SetupRedis()
	if err != nil {
		fmt.Println(err)
		cleanupRds()
		os.Exit(1)
	}

	sl = slidingLog.NewSlidingLog(client)

	exitCode := m.Run()

	cleanupRds()
	os.Exit(exitCode)
}
