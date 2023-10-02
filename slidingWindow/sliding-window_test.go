package slidingWindow_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCheckIfRequestAllowed(t *testing.T) {
	const (
		intervalToBeParsed = "10s"
		maxRequests        = 999
	)
	interval, err := time.ParseDuration(intervalToBeParsed)
	require.NoError(t, err)
	for i := 0; i < maxRequests; i++ {
		res, err := sw.CheckIfRequestAllowed("USER_IP", interval, maxRequests)
		require.NoError(t, err)
		require.True(t, res)
	}
	res, err := sw.CheckIfRequestAllowed("USER_IP", interval, maxRequests)
	require.NoError(t, err)
	require.False(t, res)
}
