package slidingLog_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCheckIfRequestAllowed(t *testing.T) {
	const (
		intervalToBeParsed = "10s"
		maxRequests        = 999
	)
	interval, err := time.ParseDuration(intervalToBeParsed)
	require.NoError(t, err)
	var res bool
	for i := 0; i < maxRequests; i++ {
		res = sl.CheckIfRequestAllowed("someUserID", uuid.NewString(), interval, maxRequests)
		require.True(t, res)
	}
	res = sl.CheckIfRequestAllowed("someUserID", uuid.NewString(), interval, maxRequests)
	require.False(t, res)
}
