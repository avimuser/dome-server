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
					return fmt.Errorf("Error running test isSetEmpty: %e", err)
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
					return fmt.Errorf("Error running test isSetEmpty: %e", err)
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

// TODO
func TestFindPair(t *testing.T) {

}
