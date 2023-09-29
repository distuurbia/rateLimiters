package slidingWindow_test

import (
	"log/slog"
	"testing"
)

func TestCheckIfRequestAllowed(t *testing.T) {
	var res bool
	for {
		res = sw.CheckIfRequestAllowed("USER_1", 1, 3)
		slog.Info("USER_1", "res", res)
		if !res {
			break
		}
		res = sw.CheckIfRequestAllowed("USER_2", 1, 3)
		slog.Info("USER_2", "res", res)
		if !res {
			break
		}
	}

}
