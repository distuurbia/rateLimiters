package tokenBucket_test

import (
	"testing"

	"github.com/distuurbia/rateLimiters/tokenBucket"
	"github.com/stretchr/testify/require"
)

func TestCheckIfRequestAllows(t *testing.T) {
	const (
		rate = 100
		requestWeight
		maxTokens = 999
	)
	tb := tokenBucket.NewTokenBucket(rate, maxTokens)

	var res bool
	for i := 0; i < 9; i++ {
		res = tb.CheckIfRequestAllowed(requestWeight)
		require.True(t, res)
	}
	res = tb.CheckIfRequestAllowed(requestWeight)
	require.False(t, res)
}
