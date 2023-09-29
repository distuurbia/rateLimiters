package fixedWindow

import (
	"sync"
	"time"
)

type FixedWindow struct {
	windowDuration      time.Duration
	lastUpdateTimestamp time.Time
	maxRequests         int64
	currentRequests     int64
	mu                  sync.Mutex
}

func NewFixedWindow(windowDuration time.Duration, maxRequests int64) *FixedWindow {
	return &FixedWindow{
		windowDuration:      windowDuration,
		maxRequests:         maxRequests,
		lastUpdateTimestamp: time.Now(),
	}
}

func (fw *FixedWindow) CheckIfRequestAllows() bool {
	now := time.Now()

	fw.mu.Lock()
	defer fw.mu.Unlock()
	if now.Sub(fw.lastUpdateTimestamp) > fw.windowDuration {
		fw.currentRequests = 0
		fw.lastUpdateTimestamp = now
	}

	if fw.currentRequests >= fw.maxRequests {
		return false
	}

	fw.currentRequests++

	return true
}
