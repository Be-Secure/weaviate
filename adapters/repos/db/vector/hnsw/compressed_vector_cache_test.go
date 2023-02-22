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

package hnsw

import (
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestCompressedVectorCacheGrowth(t *testing.T) {
	logger, _ := test.NewNullLogger()
	vectorCache := newCompressedShardedLockCache(1000000, logger)
	id := int64(10)
	assert.True(t, vectorCache.count < id)
	vectorCache.grow(uint64(id))
	assert.True(t, vectorCache.count > id)
	last := vectorCache.count
	vectorCache.grow(uint64(id))
	assert.True(t, vectorCache.count == last)
}
