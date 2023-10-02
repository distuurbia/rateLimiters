package slidingWindow

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type SlidingWindow struct {
	client *redis.Client
}

func NewSlidingWindow(client *redis.Client) *SlidingWindow {
	return &SlidingWindow{client: client}
}

func (sw *SlidingWindow) CheckIfRequestAllowed(userID string, interval time.Duration, maximumRequests int64) bool {
	now := time.Now()
	const (
		base    = 10
		bitSize = 64
	)
	intervalInSeconds := int64(interval.Seconds())
	currentWindow := strconv.FormatInt(now.Unix()/intervalInSeconds, bitSize)
	key := userID + ":" + currentWindow
	value, _ := sw.client.Get(key).Result()
	requestCountCurrentWindow, _ := strconv.ParseInt(value, base, bitSize)

	if requestCountCurrentWindow >= maximumRequests {
		return false
	}

	lastWindow := strconv.FormatInt(now.Add((-1)*interval).Unix()/intervalInSeconds, base)
	key = userID + ":" + lastWindow
	value, _ = sw.client.Get(key).Result()
	requestCountLastWindow, _ := strconv.ParseInt(value, base, bitSize)

	elapsedTimePercentage := float64(now.Unix()%intervalInSeconds) / interval.Seconds()

	if (float64(requestCountLastWindow)*(1-elapsedTimePercentage))+float64(requestCountCurrentWindow) >= float64(maximumRequests) {
		return false
	}

	sw.client.Incr(userID + ":" + currentWindow)

	return true
}
