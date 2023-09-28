package slidingLog

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type SlidingLog struct {
	client *redis.Client
}

func NewSlidingLog(client *redis.Client) *SlidingLog {
	return &SlidingLog{client: client}
}

func (sl *SlidingLog) CheckIfRequestAllowed(userID string, uniqueRequestID string, intervalInSeconds int64, maximumRequests int64) bool {
	currentTime := strconv.FormatInt(time.Now().Unix(), 10)
	lastWindowTime := strconv.FormatInt(time.Now().Unix()-intervalInSeconds, 10)
	requestCount := sl.client.ZCount(userID, lastWindowTime, currentTime).Val()
	if requestCount >= maximumRequests {
		return false
	}

	sl.client.ZAdd(userID, redis.Z{Score: float64(time.Now().Unix()), Member: uniqueRequestID})

	return true
}
