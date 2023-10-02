package slidingWindow_test

import (
	"log/slog"
	"testing"
	"time"
)

func TestCheckIfRequestAllowed(t *testing.T) {
	var res bool
	dur, _ := time.ParseDuration("10s")

	for {
		res = sw.CheckIfRequestAllowed("USER_1", dur, 3)
		slog.Info("USER_1", "res", res)
		if !res {
			break
		}
		res = sw.CheckIfRequestAllowed("USER_2", dur, 3)
		slog.Info("USER_2", "res", res)
		if !res {
			break
		}
	}

}
