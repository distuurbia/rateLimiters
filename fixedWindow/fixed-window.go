package fixedWindow

import (
	"sync"
	"time"
)

type FixedWindow struct {
	interval            time.Duration
	lastUpdateTimestamp time.Time
	maxRequests         int64
	currentRequests     int64
	mu                  sync.Mutex
}

func NewFixedWindow(interval time.Duration, maxRequests int64) *FixedWindow {
	return &FixedWindow{
		interval:            interval,
		maxRequests:         maxRequests,
		lastUpdateTimestamp: time.Now(),
	}
}

func (fw *FixedWindow) CheckIfRequestAllowed() bool {
	const firstRequest = 1
	now := time.Now()

	fw.mu.Lock()
	defer fw.mu.Unlock()
	if now.Sub(fw.lastUpdateTimestamp) > fw.interval {
		fw.currentRequests = firstRequest
		fw.lastUpdateTimestamp = now
		fw.currentRequests++
		return true
	}

	if fw.currentRequests >= fw.maxRequests {
		return false
	}

	fw.currentRequests++

	return true
}
