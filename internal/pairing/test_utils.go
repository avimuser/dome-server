package pairing

import (
	"context"
	"sync"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type redisTestCase struct {
	name          string
	setupDatabase func(*redis.Client) error
	check         func(*redis.Client) error
}

func newRedisTestConn(t *testing.T) (*redis.Client, error) {
	ctx := context.Background()
	redisC, err := testcontainers.Run(
		ctx, "redis:latest",
		testcontainers.WithExposedPorts("6379/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("6379/tcp"),
			wait.ForLog("Ready to accept connections"),
		),
	)
	if err != nil {
		return nil, err
	}
	testcontainers.CleanupContainer(t, redisC)
	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		return nil, err
	}
	conn := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})
	return conn, err
}

func runRedisTestCases(t *testing.T, cases []redisTestCase) {
	wg := sync.WaitGroup{}
	for _, testCase := range cases {
		wg.Go(func() {
			conn, err := newRedisTestConn(t)
			if err != nil {
				t.Errorf("Error making a test container for %s: %e", testCase.name, err)
			}
			defer conn.Close()

			err = testCase.setupDatabase(conn)
			if err != nil {
				t.Errorf("Error setting up database for %s: %e", testCase.name, err)
			}

			err = testCase.check(conn)
			if err != nil {
				t.Errorf("Error setting up database for %s: %e", testCase.name, err)
			}
		})
	}
	wg.Wait()
}
