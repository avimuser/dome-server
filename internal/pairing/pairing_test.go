package pairing

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestIsSetEmpty(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	testCases := []redisTestCase{
		{
			name: "Non empty set",
			setupDatabase: func(conn *redis.Client) error {
				return conn.SAdd(ctx, "test", "foo").Err()
			},
			check: func(conn *redis.Client) error {
				result, err := isSetEmpty(conn, ctx, "test")
				if err != nil {
					return fmt.Errorf("Error running test isSetEmpty: %s", err)
				}
				if result != false {
					return fmt.Errorf("Expected isSetEmpty false found true")
				}
				return nil
			},
		},
		{
			name: "Empty set",
			setupDatabase: func(conn *redis.Client) error {
				return nil
			},
			check: func(conn *redis.Client) error {
				result, err := isSetEmpty(conn, ctx, "test")
				if err != nil {
					return fmt.Errorf("Error running test isSetEmpty: %s", err)
				}
				if result != true {
					return fmt.Errorf("Expected isSetEmpty true found false")
				}
				return nil
			},
		},
	}
	runRedisTestCases(t, testCases)
}

func TestFindPair(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	testCases := []redisTestCase{
		{
			name:          "Waiting list is empty",
			setupDatabase: func(c *redis.Client) error { return nil },
			check: func(c *redis.Client) error {
				_, err := FindPair(c, ctx, "foo")
				if err != ErrWaitingListEmpty {
					return fmt.Errorf("Expected ErrWaitingListEmpty found: %s", err)
				}
				return nil
			},
		},
		{
			name: "Add two clients",
			setupDatabase: func(c *redis.Client) error {
				_, err := FindPair(c, ctx, "foo")
				if err != ErrWaitingListEmpty {
					return fmt.Errorf("Error while adding first client: %s", err)
				}
				return nil
			},
			check: func(c *redis.Client) error {
				result, err := FindPair(c, ctx, "bar")
				if err != nil {
					return fmt.Errorf("Error when finding pair: %s", err)
				}
				if result != "foo" {
					return fmt.Errorf("Expected foo found: %s", result)
				}
				return nil
			},
		},
		{
			name: "User has previous matches",
			setupDatabase: func(c *redis.Client) error {
				_, err := FindPair(c, ctx, "foo")
				if err != ErrWaitingListEmpty {
					return fmt.Errorf("Error while adding first client: %s", err)
				}
				c.SAdd(ctx, prevPrefix+"bar", "foo")
				c.SAdd(ctx, waitingList, "baz")
				return nil
			},
			check: func(c *redis.Client) error {
				result, err := FindPair(c, ctx, "bar")
				if err != nil {
					return fmt.Errorf("Error when finding pair: %s", err)
				}
				if result != "baz" {
					return fmt.Errorf("Expected baz found: %s", result)
				}
				return nil
			},
		},
	}
	runRedisTestCases(t, testCases)
}
