package slidingLog_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCheckIfRequestAllowed(t *testing.T) {
	const (
		intervalToBeParsed = "1s"
		maxRequests        = 100
	)
	ctx, cancel := context.WithCancel(context.Background())
	interval, err := time.ParseDuration(intervalToBeParsed)
	require.NoError(t, err)
	var res bool
	for i := 0; i < maxRequests; i++ {
		res = sl.CheckIfRequestAllowed(ctx, "someUserID", uuid.NewString(), interval, maxRequests)
		require.True(t, res)
	}
	res = sl.CheckIfRequestAllowed(ctx, "someUserID", uuid.NewString(), interval, maxRequests)
	require.False(t, res)
	res = sl.CheckIfRequestAllowed(ctx, "someAnotherUserID", uuid.NewString(), interval, maxRequests)
	require.True(t, res)
	cancel()
}
