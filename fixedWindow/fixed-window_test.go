package fixedWindow_test

import (
	"testing"
	"time"

	"github.com/distuurbia/rateLimiters/fixedWindow"
	"github.com/stretchr/testify/require"
)

func TestCheckIfRequestAllows(t *testing.T) {
	const (
		intervalToParsed = "10s"
		maxRequests      = 999
	)
	interval, err := time.ParseDuration(intervalToParsed)
	require.NoError(t, err)
	fw := fixedWindow.NewFixedWindow(interval, maxRequests)

	for i := 0; i < maxRequests; i++ {
		res := fw.CheckIfRequestAllowed()
		require.True(t, res)
	}
	res := fw.CheckIfRequestAllowed()
	require.False(t, res)
}
