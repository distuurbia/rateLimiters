package fixedWindow

import "time"

type FixedWindow struct {
	windowTime          time.Duration
	lastUpdateTimestamp time.Time
	maxRequests         int64
	currentRequests     int64
}

func NewFixedWindow(windowTime time.Duration, maxRequests int64) *FixedWindow {
	return &FixedWindow{windowTime: windowTime, maxRequests: maxRequests, lastUpdateTimestamp: time.Now()}
}

func (fw *FixedWindow) CheckIfRequestAllows() bool {
	now := time.Now()

	if now.Sub(fw.lastUpdateTimestamp) > fw.windowTime {
		fw.currentRequests = 0
		fw.lastUpdateTimestamp = now
	}

	if fw.currentRequests >= fw.maxRequests {
		return false
	}

	fw.currentRequests++

	return true
}
