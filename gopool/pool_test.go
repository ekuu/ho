package gopool

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_pool_Wait(t *testing.T) {
	var i int
	Go(func(ctx context.Context) {
		i = 1
	})
	i = 2
	Wait()
	require.Equal(t, i, 1)
}
