package slidingLog_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCheckIfRequestAllowed(t *testing.T) {
	var res bool
	for i := 0; i < 1000; i++ {
		res = sl.CheckIfRequestAllowed("someUserID", uuid.NewString(), 1, 999)
	}
	require.False(t, res)
}
