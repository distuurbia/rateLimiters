package slidingLog

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type SlidingLog struct {
	client *redis.Client
}

func NewSlidingLog(client *redis.Client) *SlidingLog {
	return &SlidingLog{client: client}
}

func (sl *SlidingLog) CheckIfRequestAllowed(ctx context.Context, userID string, uniqueRequestID string, interval time.Duration, maximumRequests int64) bool {
	const base = 10
	currentTime := strconv.FormatInt(time.Now().Unix(), base)
	lastWindowTime := strconv.FormatInt(time.Now().Add((-1)*interval).Unix(), base)
	requestCount := sl.client.ZCount(ctx, userID, lastWindowTime, currentTime).Val()
	if requestCount >= maximumRequests {
		return false
	}

	sl.client.ZAdd(ctx, userID, &redis.Z{Score: float64(time.Now().Unix()), Member: uniqueRequestID})
	sl.client.ZRemRangeByScore(ctx, userID, "", lastWindowTime)
	return true
}
