package pairing

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	ErrWaitingListEmpty     = fmt.Errorf("Waiting list is empty")
	ErrFailedWaitingListAdd = fmt.Errorf("Failed to add to waiting list")
)

const (
	waitingList = "waitingList"
	prevPrefix  = "prev:"
)

func isSetEmpty(conn *redis.Client, ctx context.Context, set string) (bool, error) {
	setCard, err := conn.SCard(ctx, set).Result()
	if err != nil {
		return false, err
	}
	if setCard <= 0 {
		return true, nil
	}
	return false, nil
}

func getRandomPair(conn *redis.Client, ctx context.Context) (string, error) {
	return conn.SRandMember(ctx, waitingList).Result()
}

func FindPair(conn *redis.Client, ctx context.Context, client string) (string, error) {
	isWLEmpty, err := isSetEmpty(conn, ctx, waitingList)
	if err != nil {
		return "", fmt.Errorf("Check waiting list empty: %e", err)
	}
	if isWLEmpty {
		err := conn.SAdd(ctx, waitingList, client).Err()
		if err != nil {
			return "", ErrFailedWaitingListAdd
		}
		return "", ErrWaitingListEmpty
	}

	isPrevEmpty, err := isSetEmpty(conn, ctx, prevPrefix+client)
	if err != nil {
		return "", fmt.Errorf("Check previous matches: %e", err)
	}
	if isPrevEmpty {
		randomPair, err := getRandomPair(conn, ctx)
		if err != nil {
			return "", fmt.Errorf("Check previous matches list empty: %e", err)
		}
		return randomPair, nil
	}

	set, err := conn.SDiff(ctx, waitingList, prevPrefix+client).Result()
	if err != nil {
		return "", fmt.Errorf("Get difference b/w waiting list and previous clients: %e", err)
	}
	fmt.Println(set)
	if len(set) <= 0 {
		randomPair, err := getRandomPair(conn, ctx)
		if err != nil {
			return "", fmt.Errorf("Check previous matches list empty: %e", err)
		}
		return randomPair, nil
	}
	return set[0], nil
}
