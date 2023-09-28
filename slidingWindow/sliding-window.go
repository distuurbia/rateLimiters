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

func (sw *SlidingWindow) CheckIfRequestAllowed(userID string, intervalInSeconds int64, maximumRequests int64) bool {
	now := time.Now().Unix()

	currentWindow := strconv.FormatInt(now/intervalInSeconds, 10)
	key := userID + ":" + currentWindow
	value, _ := sw.client.Get(key).Result()
	requestCountCurrentWindow, _ := strconv.ParseInt(value, 10, 64)
	if requestCountCurrentWindow >= maximumRequests {
		return false
	}
	lastWindow := strconv.FormatInt((now-intervalInSeconds)/intervalInSeconds, 10)
	key = userID + ":" + lastWindow
	value, _ = sw.client.Get(key).Result()
	requestCountLastWindow, _ := strconv.ParseInt(value, 10, 64)

	elapsedTimePercentage := float64(now%intervalInSeconds) / float64(intervalInSeconds)

	if (float64(requestCountLastWindow)*(1-elapsedTimePercentage))+float64(requestCountCurrentWindow) >= float64(maximumRequests) {
		return false
	}

	sw.client.Incr(userID + ":" + currentWindow)

	return true
}
