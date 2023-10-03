package slidingWindow

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type SlidingWindow struct {
	client *redis.Client
}

func NewSlidingWindow(client *redis.Client) *SlidingWindow {
	return &SlidingWindow{client: client}
}

func (sw *SlidingWindow) CheckIfRequestAllowed(ctx context.Context, userID string, interval time.Duration, maximumRequests int64) (bool, error) {
	now := time.Now()
	const (
		base    = 10
		bitSize = 64
	)
	intervalInSeconds := int64(interval.Seconds())
	currentWindow := strconv.FormatInt(now.Unix()/intervalInSeconds, base)
	key := userID + ":" + currentWindow

	value, err := sw.client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return false, fmt.Errorf("slidingWindow-sw.client.Get(key).Result-err: %w", err)
	}

	var requestCountCurrentWindow int64
	if value == "" {
		requestCountCurrentWindow = 0
	} else {
		requestCountCurrentWindow, err = strconv.ParseInt(value, base, bitSize)
		if err != nil {
			return false, fmt.Errorf("slidingWindow-strconv.ParseInt-err: %w", err)
		}
	}
	if requestCountCurrentWindow >= maximumRequests {
		return false, nil
	}

	lastWindow := strconv.FormatInt(now.Add((-1)*interval).Unix()/intervalInSeconds, base)
	key = userID + ":" + lastWindow
	value, err = sw.client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return false, fmt.Errorf("slidingWindow-sw.client.Get(key).Result-err: %w", err)
	}

	var requestCountLastWindow int64
	if value == "" {
		requestCountCurrentWindow = 0
	} else {
		requestCountLastWindow, err = strconv.ParseInt(value, base, bitSize)
		if err != nil {
			return false, fmt.Errorf("slidingWindow-strconv.ParseInt-err: %w", err)
		}
	}

	elapsedTimePercentage := float64(now.Unix()%intervalInSeconds) / interval.Seconds()

	if (float64(requestCountLastWindow)*(1-elapsedTimePercentage))+float64(requestCountCurrentWindow) >= float64(maximumRequests) {
		return false, nil
	}

	sw.client.Incr(ctx, userID+":"+currentWindow)
	sw.client.Expire(ctx, userID+":"+currentWindow, interval)

	return true, nil
}
