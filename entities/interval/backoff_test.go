//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2023 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package interval

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBackoffInterval(t *testing.T) {
	t.Run("with default backoffs", func(t *testing.T) {
		boff := NewBackoffTimer()

		assert.Equal(t, boff.backoffs, defaultBackoffs)
		assert.Zero(t, boff.backoffLevel)
		assert.Zero(t, boff.lastInterval)
		assert.Equal(t, time.Duration(0), boff.getWarningInterval())
		assert.True(t, boff.IntervalElapsed())

		i := 1
		for ; i < len(defaultBackoffs); i++ {
			boff.IncreaseInterval()
			assert.False(t, boff.IntervalElapsed())
			assert.Equal(t, i, boff.backoffLevel)
			assert.Equal(t, defaultBackoffs[i], boff.getWarningInterval())
		}

		boff.IncreaseInterval()
		assert.False(t, boff.IntervalElapsed())
		assert.Equal(t, i, boff.backoffLevel)
		assert.Equal(t, 24*time.Hour, boff.getWarningInterval())
	})

	t.Run("with custom backoffs", func(t *testing.T) {
		var (
			durations = []time.Duration{time.Second, time.Nanosecond, time.Millisecond}
			sorted    = make([]time.Duration, len(durations))
		)

		copy(sorted, durations)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i] < sorted[j]
		})

		boff := NewBackoffTimer(durations...)
		assert.Equal(t, boff.backoffs, sorted)
		assert.True(t, boff.IntervalElapsed())
		assert.Equal(t, sorted[0], boff.getWarningInterval())

		boff.IncreaseInterval()
		time.Sleep(time.Millisecond)
		assert.True(t, boff.IntervalElapsed())
		assert.Equal(t, sorted[1], boff.getWarningInterval())

		boff.IncreaseInterval()
		assert.False(t, boff.IntervalElapsed())
		time.Sleep(time.Second)
		assert.True(t, boff.IntervalElapsed())
		assert.Equal(t, sorted[2], boff.getWarningInterval())

		boff.IncreaseInterval()
		assert.False(t, boff.IntervalElapsed())
		assert.False(t, boff.IntervalElapsed())
		assert.Equal(t, 24*time.Hour, boff.getWarningInterval())
	})
}
